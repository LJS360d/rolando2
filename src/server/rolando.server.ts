import { ChannelType, Client, TextChannel } from 'discord.js';
import express, { Request, Response } from 'express';
import { Fonzi2Server, Fonzi2ServerData, Logger } from 'fonzi2';
import { resolve } from 'path';
import { MarkovChainAnalyzer } from '../domain/model/chain.analyzer';
import { ChainsService } from '../domain/services/chains.service';

export class RolandoServer extends Fonzi2Server {
	constructor(
		client: Client,
		data: Fonzi2ServerData,
		private chainsService: ChainsService
	) {
		super(client, data);
		this.app.use(express.static(resolve('public')));
		this.app.set('views', [this.app.get('views'), resolve('views')]);
	}

	override async start() {
		this.app.get('/chain', this.guildChain.bind(this));
		super.start();
	}

	override async dashboard(req: Request, res: Response) {
		const userInfo = req.session!['userInfo'];
		if (!userInfo) {
			res.redirect('/unauthorized');
			return;
		}

		const invites = {};

		const invitePromises = this.client.guilds.cache.map(async (guild) => {
			const channel = guild.channels.cache.find(
				(channel) => channel && channel.type === ChannelType.GuildText
			) as TextChannel | undefined;

			if (channel) {
				try {
					const invite = await channel.createInvite();
					invites[guild.id] = `https://discord.gg/${invite.code}`;
				} catch (error) {
					Logger.warn(`Missing invite permissions in ${guild.name}`);
					invites[guild.id] = '';
				}
				return;
			}
		});

		await Promise.all(invitePromises);

		const props = {
			client: this.client,
			guilds: this.client.guilds.cache,
			startTime: this.startTime,
			version: this.data.version,
			inviteLink: this.data.inviteLink,
			userInfo,
			//? Rolando specific
			chains: this.chainsService.chains,
			invites,
			analyzer: MarkovChainAnalyzer,
		};

		res.render('dashboard', props);
	}

	private async guildChain(req: Request, res: Response) {
		if (!req.session!['userInfo']) {
			res.redirect('/unauthorized');
			return;
		}
		const { guildId } = req.query;
		if (typeof guildId === 'string') {
			const chain = this.chainsService.chains.get(guildId) ?? {
				code: 404,
				message: `chain ${guildId} not found`,
			};
			res.status(!!chain ? 200 : 404).json(chain);
		}
	}
}
