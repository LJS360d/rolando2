import { MediaStorage } from './media.storage';

export class MarkovChain {
	public mediaStorage: MediaStorage;
	state: MarkovState;
	messageCounter = 0;

	constructor(
		public id: string,
		public replyRate = 10,
		public pings = true,
		messages: string[] = []
	) {
		this.mediaStorage = new MediaStorage(this.id);
		this.state = {};
		this.provideData(messages);
	}

	provideData(messages: string[]): void {
		messages.forEach((message) => this.updateState(message));
	}

	updateState(message: string): void {
		if (message.startsWith('https://')) {
			this.mediaStorage.addMedia(message);
			return;
		}
		this.messageCounter++;
		const tokens = this.tokenize(message);
		for (let i = 0; i < tokens.length - 1; i++) {
			const currentWord = tokens[i];
			const nextWord = tokens[i + 1];
			this.state[currentWord] ??= {};
			this.state[currentWord][nextWord] ??= 0;
			this.state[currentWord][nextWord]++;
		}
	}

	generateText(startWord: string, length: number): string {
		let currentWord = startWord;
		let generatedText = currentWord;
		for (let i = 0; i < length; i++) {
			const nextWords = this.state[currentWord];
			if (!nextWords) {
				break;
			}
			const nextWordArray = Object.keys(nextWords);
			const nextWordWeights = Object.values(nextWords);
			// ? Laplace smoothing
			const smoothedWeights = nextWordWeights.map(
				(weight) => (weight + 1) / (this.messageCounter + nextWordArray.length)
			);
			currentWord = this.stochasticChoice(nextWordArray, smoothedWeights);
			generatedText += ` ${currentWord}`;
		}
		if (!this.pings) {
			generatedText = generatedText.replace(/<(@&?\w+)>/g, '$1');
		}
		return generatedText;
	}

	talk(length: number): string {
		const keys = Object.keys(this.state);
		const randomIndex = Math.floor(Math.random() * keys.length);
		const starterWord = keys[randomIndex];
		return this.generateText(starterWord, length).trim();
	}

	delete(message: string) {
		if (message.startsWith('https://')) {
			this.mediaStorage.removeMedia(message);
		}
		const tokens = this.tokenize(message);
		for (let i = 0; i < tokens.length - 1; i++) {
			const currentWord = tokens[i];
			const nextWord = tokens[i + 1];
			if (this.state[currentWord]?.[nextWord]) {
				delete this.state[currentWord][nextWord];
			}
		}
	}

	private tokenize(text: string) {
		const tokens = text.split(/\s+/);
		return tokens.filter((token) => token.length > 0);
	}

	private stochasticChoice(options: string[], weights: number[]): string {
		const totalWeight = weights.reduce((a, b) => a + b, 0);
		const randomWeight = Math.random() * totalWeight;
		let weightSum = 0;
		for (let i = 0; i < options.length; i++) {
			weightSum += weights[i];
			if (randomWeight <= weightSum) {
				return options[i];
			}
		}
		return options[options.length - 1];
	}
}
