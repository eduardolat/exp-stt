/**
 * Formats a timestamp into a relative time string (e.g., "5 minutes ago").
 * Similar to date-fns formatDistanceToNow({ addSuffix: true }).
 */
export function formatRelativeTime(dateOrString: Date | string): string {
	const date = new Date(dateOrString);
	const now = new Date();
	const diffMs = now.getTime() - date.getTime();
	const diffSec = Math.round(diffMs / 1000);

	const rtf = new Intl.RelativeTimeFormat('en', { numeric: 'auto' });

	if (diffSec < 60) {
		return 'less than a minute ago';
	}
	const diffMin = Math.round(diffSec / 60);
	if (diffMin < 60) {
		return rtf.format(-diffMin, 'minute');
	}
	const diffHour = Math.round(diffMin / 60);
	if (diffHour < 24) {
		return rtf.format(-diffHour, 'hour');
	}
	const diffDay = Math.round(diffHour / 24);
	if (diffDay < 7) {
		return rtf.format(-diffDay, 'day');
	}
	const diffMonth = Math.round(diffDay / 30);
	if (diffMonth < 12) {
		return rtf.format(-diffMonth, 'month');
	}
	const diffYear = Math.round(diffMonth / 12);
	return rtf.format(-diffYear, 'year');
}
