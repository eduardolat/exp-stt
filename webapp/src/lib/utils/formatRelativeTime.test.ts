import { describe, it, expect, vi, afterEach, beforeEach } from 'vitest';
import { formatRelativeTime } from './formatRelativeTime';

describe('formatRelativeTime', () => {
	beforeEach(() => {
		vi.useFakeTimers();
	});

	afterEach(() => {
		vi.useRealTimers();
	});

	it('returns "just now" for times less than a minute ago', () => {
		const now = new Date('2023-01-01T12:00:00Z');
		vi.setSystemTime(now);
		expect(formatRelativeTime(new Date('2023-01-01T12:00:00Z'))).toBe('just now');
		expect(formatRelativeTime(new Date('2023-01-01T11:59:30Z'))).toBe('just now');
	});

	it('returns minutes for times less than an hour ago', () => {
		const now = new Date('2023-01-01T12:00:00Z');
		vi.setSystemTime(now);
		expect(formatRelativeTime(new Date('2023-01-01T11:59:00Z'))).toBe('1 minute ago');
		expect(formatRelativeTime(new Date('2023-01-01T11:50:00Z'))).toBe('10 minutes ago');
	});

	it('returns hours for times less than a day ago', () => {
		const now = new Date('2023-01-01T12:00:00Z');
		vi.setSystemTime(now);
		expect(formatRelativeTime(new Date('2023-01-01T11:00:00Z'))).toBe('1 hour ago');
		expect(formatRelativeTime(new Date('2023-01-01T02:00:00Z'))).toBe('10 hours ago');
	});

	it('returns days for times less than a month ago', () => {
		const now = new Date('2023-01-01T12:00:00Z');
		vi.setSystemTime(now);
		expect(formatRelativeTime(new Date('2022-12-31T12:00:00Z'))).toBe('yesterday');
		expect(formatRelativeTime(new Date('2022-12-20T12:00:00Z'))).toBe('12 days ago');
	});

	it('returns months for times less than a year ago', () => {
		const now = new Date('2023-01-01T12:00:00Z');
		vi.setSystemTime(now);
		expect(formatRelativeTime(new Date('2022-12-01T12:00:00Z'))).toBe('last month');
		expect(formatRelativeTime(new Date('2022-06-01T12:00:00Z'))).toBe('7 months ago');
	});

	it('returns years for times more than a year ago', () => {
		const now = new Date('2023-01-01T12:00:00Z');
		vi.setSystemTime(now);
		expect(formatRelativeTime(new Date('2022-01-01T12:00:00Z'))).toBe('last year');
		expect(formatRelativeTime(new Date('2010-01-01T12:00:00Z'))).toBe('13 years ago');
	});

	it('handles string input', () => {
		const now = new Date('2023-01-01T12:00:00Z');
		vi.setSystemTime(now);
		expect(formatRelativeTime('2023-01-01T11:50:00Z')).toBe('10 minutes ago');
	});

	it('handles number input', () => {
		const now = new Date('2023-01-01T12:00:00Z');
		vi.setSystemTime(now);
		expect(formatRelativeTime(new Date('2023-01-01T11:50:00Z').getTime())).toBe('10 minutes ago');
	});
});
