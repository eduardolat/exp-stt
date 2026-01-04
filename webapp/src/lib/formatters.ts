export function formatTimeAgo(date: Date | string | number): string {
	const d = new Date(date);
	const now = new Date();
	const diffInSeconds = (d.getTime() - now.getTime()) / 1000;
	const rtf = new Intl.RelativeTimeFormat('en', { numeric: 'always' });

	// Use absolute value for comparisons, but keep sign for formatting
	const absDiff = Math.abs(diffInSeconds);

	if (absDiff < 60) {
		return 'just now';
	} else if (absDiff < 3600) {
		return rtf.format(Math.round(diffInSeconds / 60), 'minute');
	} else if (absDiff < 86400) {
		return rtf.format(Math.round(diffInSeconds / 3600), 'hour');
	} else if (absDiff < 2592000) {
		// 30 days
		return rtf.format(Math.round(diffInSeconds / 86400), 'day');
	} else if (absDiff < 31536000) {
		// 365 days
		return rtf.format(Math.round(diffInSeconds / 2592000), 'month');
	} else {
		return rtf.format(Math.round(diffInSeconds / 31536000), 'year');
	}
}
