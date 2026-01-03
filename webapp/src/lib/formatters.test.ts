import { describe, it, expect } from 'vitest';
import { formatTimeAgo } from './formatters';

describe('formatTimeAgo', () => {
    const now = new Date();

    it('returns "just now" for less than 60 seconds', () => {
        const d = new Date(now.getTime() - 30 * 1000);
        expect(formatTimeAgo(d)).toBe('just now');
    });

    it('returns minutes ago', () => {
        const d = new Date(now.getTime() - 5 * 60 * 1000);
        // formatTimeAgo uses relative time format, which might be "5 minutes ago"
        expect(formatTimeAgo(d)).toBe('5 minutes ago');
    });

    it('returns hours ago', () => {
        const d = new Date(now.getTime() - 2 * 60 * 60 * 1000);
        expect(formatTimeAgo(d)).toBe('2 hours ago');
    });

    it('returns days ago', () => {
        const d = new Date(now.getTime() - 3 * 24 * 60 * 60 * 1000);
        expect(formatTimeAgo(d)).toBe('3 days ago');
    });

    it('returns months ago', () => {
        const d = new Date(now.getTime() - 60 * 24 * 60 * 60 * 1000); // 60 days ~ 2 months
        expect(formatTimeAgo(d)).toBe('2 months ago');
    });

    it('returns years ago', () => {
        const d = new Date(now.getTime() - 400 * 24 * 60 * 60 * 1000); // > 365 days
        expect(formatTimeAgo(d)).toBe('1 year ago');
    });

    it('handles string input', () => {
        const d = new Date(now.getTime() - 30 * 1000);
        expect(formatTimeAgo(d.toISOString())).toBe('just now');
    });
});
