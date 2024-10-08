import { ChannelType, type Client, type TextChannel } from 'discord.js';
import express, { type Request, type Response } from 'express';
import {
	Fonzi2Server,
	Logger,
	type ServerController,
	getRegisteredCommands,
} from 'fonzi2';
import { resolve } from 'node:path';
import { MarkovChainAnalyzer } from '../domain/model/chain.analyzer';
import type { ChainsService } from '../domain/services/chains.service';
import { env } from '../env';
import { render } from './render';

export class RolandoServer extends Fonzi2Server {
	constructor(
		client: Client<true>,
		private chainsService: ChainsService,
		controllers: ServerController[]
	) {
		super(client, controllers);
		this.app.use(express.urlencoded({ extended: true }));
		this.app.use(express.static(resolve('public')));
		this.app.set('views', [this.app.get('views'), resolve('views')]);
	}

	override async start() {
		this.app.get('/chain', this.guildChain.bind(this));
		this.app.get('/data', this.guildMessages.bind(this));
		this.app.get('/invite', this.getGuildInvite.bind(this));
		this.app.get('/home', this.home.bind(this));
		this.app.get('/tos', this.tos.bind(this));
		this.app.get('/privacy', this.privacy.bind(this));
		this.app.get('/broadcast', this.broadcast.bind(this));
		this.app.post('/broadcast', this.sendBroadcast.bind(this));
		this.app.get('/chains/memUsage', this.memUsage.bind(this));
		super.start();
	}

	override async dashboard(req: Request, res: Response) {
		const userInfo = this.getSessionUserInfo(req);
		if (!userInfo) {
			res.redirect('/');
			return;
		}
		if (userInfo.role !== 'owner') {
			res.redirect('/home');
			return;
		}

		const props = {
			client: this.client,
			guilds: this.client.guilds.cache,
			startTime: this.startTime,
			inviteLink: env.INVITE_LINK,
			//? Rolando specific
			chains: this.chainsService.chains,
			MarkovChainAnalyzer,
		};

		const options = {
			version: env.VERSION,
			userInfo,
		};
		render(res, 'pages/backoffice/dashboard', props, options);
		return;
	}

	async broadcast(req: Request, res: Response) {
		const userInfo = this.getSessionUserInfo(req);
		if (!userInfo) {
			res.redirect('/');
			return;
		}
		if (userInfo.role !== 'owner') {
			res.redirect('/home');
			return;
		}

		const props = {
			client: this.client,
			guilds: this.client.guilds.cache,
		};

		const options = {
			version: env.VERSION,
			userInfo,
		};
		render(res, 'pages/backoffice/broadcast', props, options);
		return;
	}

	async home(req: Request, res: Response) {
		const props = {
			client: this.client,
			guilds: this.client.guilds.cache.size,
			inviteLink: env.INVITE_LINK,
			commands: getRegisteredCommands(),
		};

		const options = {
			version: env.VERSION,
			userInfo: this.getSessionUserInfo(req),
		};
		render(res, 'pages/frontoffice/home', props, options);
		return;
	}

	async tos(req: Request, res: Response) {
		const props = {
			client: this.client,
		};

		const options = {
			version: env.VERSION,
			userInfo: this.getSessionUserInfo(req),
		};
		render(res, 'pages/frontoffice/tos', props, options);
		return;
	}

	async privacy(req: Request, res: Response) {
		const props = {
			client: this.client,
		};

		const options = {
			version: env.VERSION,
			userInfo: this.getSessionUserInfo(req),
		};
		render(res, 'pages/frontoffice/privacy', props, options);
		return;
	}

	private async getGuildInvite(
		req: Request<any, any, any, { guildId: string }>,
		res: Response
	) {
		const guild = this.client.guilds.cache.get(req.query.guildId)!;
		const channel = guild.channels.cache.find(
			(channel) => channel && channel.type === ChannelType.GuildText
		) as TextChannel | undefined;
		if (!channel) {
			res.send(
				'<a><i style="color: red;" class="fa-solid fa-door-closed"></i></a>'
			);
			return;
		}
		try {
			const invite = await channel.createInvite();
			res.send(
				`<a href="https://discord.gg/${invite.code}"><i style="color: lime;" class="fa-solid fa-door-open"></i></a>`
			);
		} catch (error) {
			res.send(
				'<a><i style="color: red;" class="fa-solid fa-door-closed"></i></a>'
			);
			return;
		}
	}

	private async guildChain(
		req: Request<any, any, any, { guildId: string }>,
		res: Response
	) {
		if (!this.getSessionUserInfo(req)) {
			res.sendStatus(401);
			return;
		}
		const { guildId } = req.query;
		const chain = this.chainsService.chains.get(guildId) ?? {
			code: 404,
			message: `chain ${guildId} not found`,
		};
		res.status(chain ? 200 : 404).json(chain);
	}

	private async guildMessages(
		req: Request<any, any, any, { guildId: string }>,
		res: Response
	) {
		if (!this.getSessionUserInfo(req)) {
			res.sendStatus(401);
			return;
		}
		const { guildId } = req.query;
		const messages = {
			messages: this.chainsService.getChainMessages(guildId),
		} ?? {
			code: 404,
			message: `chain ${guildId} not found`,
		};
		res.status(messages ? 200 : 404).json(messages);
	}

	private async memUsage(req: Request, res: Response) {
		const chainsMemUsage = this.chainsService.getChainsMemUsage();
		res
			.status(200)
			.send(
				`<span class="text-sm">Chains memory usage:<b>${chainsMemUsage}</b></span>`
			);
	}

	private async sendBroadcast(req: Request, res: Response) {
		const userInfo = this.getSessionUserInfo(req);
		if (!userInfo || userInfo.role !== 'owner') {
			res.sendStatus(403);
			return;
		}
		const msg = req.body.msg as string;
		if (!msg) {
			res.sendStatus(400);
			return;
		}
		const guildIds = Object.entries(req.body)
			.filter(([key, value]) => value === 'on' && key !== 'msg')
			.map(([key, _]) => key);
		Logger.info(`Sending Broadcast message: ${msg}`);
		this.client.guilds.cache.forEach((guild) => {
			if (guildIds.includes(guild.id)) {
				guild.systemChannel?.send(msg);
			}
		});
		res.sendStatus(200);
	}
}
