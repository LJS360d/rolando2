import { getValidUrl, isGifUrl, isImageUrl, isVideoUrl } from '../../utils/url.utils';

export class MediaStorage {
	gifs: Set<string>;
	images: Set<string>;
	videos: Set<string>;
	constructor(public chainId: string) {
		this.gifs = new Set<string>();
		this.images = new Set<string>();
		this.videos = new Set<string>();
	}

	addMedia(url: string) {
		if (isGifUrl(url)) this.gifs.add(url);
		else if (isVideoUrl(url)) this.videos.add(url);
		else if (isImageUrl(url)) this.images.add(url);
		return;
	}

	async getMedia(type: 'gif' | 'video' | 'image'): Promise<string> {
		switch (type) {
			case 'gif':
				return getValidUrl(this.gifs, type);
			case 'video':
				return getValidUrl(this.videos, type);
			case 'image':
				return getValidUrl(this.images, type);
		}
	}

	removeMedia(url: string): void {
		this.gifs.delete(url);
		this.videos.delete(url);
		this.images.delete(url);
	}
}