/**
 * ADD MORE ICONS HERE IF NEEDED
 */
import { Sparkles, GraduationCap, MessageCircle, FileText, Languages, Code } from '@lucide/svelte';

const ICON_MAP: Record<string, Component> = {
	Sparkles,
	GraduationCap,
	MessageCircle,
	FileText,
	Languages,
	Code
};

/**
 * REST OF THE LOGIC
 */

import type { Component } from 'svelte';

export interface PromptPreset {
	id: string;
	name: string;
	description: string;
	icon: Component;
	body: string;
}

interface PresetFrontmatter {
	name: string;
	description: string;
	icon: string;
}

function parseFrontmatter(content: string): { frontmatter: PresetFrontmatter; body: string } {
	const match = content.match(/^---\n([\s\S]*?)\n---\n([\s\S]*)$/);
	if (!match) {
		throw new Error('Invalid frontmatter format');
	}

	const [, frontmatterStr, body] = match;
	const frontmatter: Partial<PresetFrontmatter> = {};

	for (const line of frontmatterStr.split('\n')) {
		const colonIndex = line.indexOf(':');
		if (colonIndex === -1) continue;

		const key = line.slice(0, colonIndex).trim();
		const value = line.slice(colonIndex + 1).trim();

		if (key === 'name' || key === 'description' || key === 'icon') {
			frontmatter[key] = value;
		}
	}

	return {
		frontmatter: frontmatter as PresetFrontmatter,
		body: body.trim()
	};
}

function getIdFromPath(path: string): string {
	const match = path.match(/\/([^/]+)\.md$/);
	return match ? match[1] : path;
}

const presetModules = import.meta.glob('./*.md', { query: '?raw', eager: true, import: 'default' });

export const promptPresets: PromptPreset[] = Object.entries(presetModules).map(
	([path, content]) => {
		const { frontmatter, body } = parseFrontmatter(content as string);
		return {
			id: getIdFromPath(path),
			name: frontmatter.name,
			description: frontmatter.description,
			icon: ICON_MAP[frontmatter.icon] ?? FileText,
			body
		};
	}
);
