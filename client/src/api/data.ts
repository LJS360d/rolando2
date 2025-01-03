import { useQuery } from "@tanstack/vue-query";

export interface GuildData {
  guild: GuildMeta;
  messages: string[];
}

export interface GuildMeta {
  icon: string;
  id: string;
  members: number;
  name: string;
}


export function useGetGuildData(token: string, guildId: string) {
  return useQuery({
    queryKey: ["/bot/resources"],
    queryFn: async () => {
      const response = await fetch(`/api/data/${guildId}`, {
        headers: {
          Authorization: token
        }
      });
      if (!response.ok) throw new Error("Failed to fetch guild data");
      return response.json() as Promise<GuildData>;
    },
  });
}