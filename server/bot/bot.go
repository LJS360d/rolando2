package bot

import (
	context "context"
	"rolando/cmd/log"
	"rolando/config"
	sync "sync"

	"github.com/bwmarrin/discordgo"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type BotServerImpl struct {
	BotServer
	ds *discordgo.Session
}

func NewBotServer(ds *discordgo.Session) *BotServerImpl {
	return &BotServerImpl{
		ds: ds,
	}
}

func (s *BotServerImpl) GetBotUser(ctx context.Context, req *emptypb.Empty) (*BotUser, error) {
	botUser := s.ds.State.User
	commands, err := s.getSlashCommands()
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	return &BotUser{
		Id:            botUser.ID,
		Username:      botUser.Username,
		GlobalName:    botUser.Username + "#" + botUser.Discriminator,
		AvatarUrl:     botUser.AvatarURL(""),
		Discriminator: botUser.Discriminator,
		Verified:      botUser.Verified,
		AccentColor:   int32(botUser.AccentColor),
		InviteUrl:     config.InviteUrl,
		SlashCommands: commands,
	}, nil
}

func (s *BotServerImpl) GetBotGuilds(req *emptypb.Empty, stream grpc.ServerStreamingServer[BotGuild]) error {
	defer stream.Context().Done()
	guilds, err := s.ds.UserGuilds(200, "", "", true)
	if err != nil {
		stream.Context().Err()
		return status.Error(codes.FailedPrecondition, err.Error())
	}
	for _, guild := range guilds {
		guildChannel, err := s.ds.GuildChannels(guild.ID)
		if err != nil {
			stream.Context().Err()
			return status.Error(codes.FailedPrecondition, err.Error())
		}
		botGuild := &BotGuild{
			Id:   guild.ID,
			Name: guild.Name,
		}
		for _, channel := range guildChannel {
			botGuild.Channels = append(botGuild.Channels, &BotGuildChannel{
				Id:   channel.ID,
				Name: channel.Name,
				Text: channel.Type == discordgo.ChannelTypeGuildText,
			})
		}
		if err := stream.Send(botGuild); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
	// Sends EOF
	return nil
}

func (s *BotServerImpl) BroadcastMessage(ctx context.Context, req *BroadcastMessageRequest) (*emptypb.Empty, error) {
	var wg sync.WaitGroup
	errCh := make(chan error, len(req.Guilds))

	for _, g := range req.Guilds {
		wg.Add(1)

		go func(g *BroadcastGuildRequest) {
			defer wg.Done()

			channelId := g.GetChannelId()
			if channelId == "" {
				guild, err := s.ds.Guild(g.GetId())
				if err != nil {
					errCh <- status.Error(codes.FailedPrecondition, err.Error())
					return
				}
				channelId = guild.SystemChannelID
			}
			log.Log.Infof("Broadcasting message in guild %s", g.GetId())
			_, err := s.ds.ChannelMessageSend(channelId, req.Content)
			if err != nil {
				log.Log.Errorf("could not send message in guild %s, channel: %s: %v", g.GetId(), channelId, err)
				errCh <- status.Error(codes.FailedPrecondition, err.Error())
			}
		}(g)
	}

	wg.Wait()
	close(errCh)

	// Collect errors, if any
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	return &emptypb.Empty{}, nil
}

// -------------- Helpers --------------

func (s *BotServerImpl) getSlashCommands() ([]*SlashCommand, error) {
	commands, err := s.ds.ApplicationCommands(s.ds.State.User.ID, "")
	if err != nil {
		return nil, err
	}
	slashCommands := make([]*SlashCommand, len(commands))
	for i, command := range commands {
		slashCommands[i] = &SlashCommand{
			Name:        command.Name,
			Description: command.Description,
		}
	}
	return slashCommands, nil
}
