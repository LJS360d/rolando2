import { useQuery } from "@tanstack/vue-query";

export interface BotUser {
  accent_color: number;
  avatar_url: string;
  discriminator: string;
  global_name: string;
  id: string;
  invite_url: string;
  slash_commands: SlashCommand[];
  username: string;
  verified: boolean;
  guilds: number;
  startup_timestamp_unix: number;
  mem_usage_peak: bigint;
  mem_usage_max: bigint;
}

export interface SlashCommand {
  id: string;
  application_id: string;
  version: string;
  type: number;
  name: string;
  dm_permission: boolean;
  nsfw: boolean;
  description: string;
  options: Option[] | null;
}

export interface Option {
  type: number;
  name: string;
  description: string;
  channel_types: null;
  required: boolean;
  options: null;
  autocomplete: boolean;
  choices: null;
}

export function useGetBotUser() {
  return useQuery({
    queryKey: ["/bot/user"],
    queryFn: async () => {
      const response = await fetch(`/api/bot/user`);
      if (!response.ok) throw new Error("Failed to fetch bot user");
      return response.json() as Promise<BotUser>;
    },
  });
}

export interface BotGuild {
  id: string;
  name: string;
  icon: string;
  owner: boolean;
  permissions: string;
  features: string[];
  approximate_member_count: number;
  approximate_presence_count: number;
}

export function useGetBotGuilds(token: string) {
  return useQuery({
    queryKey: ["/bot/guilds"],
    queryFn: async () => {
      const response = await fetch(`/api/bot/guilds`, {
        headers: {
          Authorization: token
        }
      });
      if (!response.ok) throw new Error("Failed to fetch bot guilds");
      return response.json() as Promise<BotGuild[]>;
    },
  });
}