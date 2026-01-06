
export function formatRelativeTime(date: Date | string | number): string {
	const d = new Date(date);
	const now = new Date();
	const diffMs = d.getTime() - now.getTime();
	const diffSecs = Math.round(diffMs / 1000);

	const rtf = new Intl.RelativeTimeFormat('en', { numeric: 'auto' });

	// Helper to handle "just now"
	if (Math.abs(diffSecs) < 10) {
		return 'just now';
	}

	if (Math.abs(diffSecs) < 60) {
		return rtf.format(diffSecs, 'second');
	}

	const diffMins = Math.round(diffSecs / 60);
	if (Math.abs(diffMins) < 60) {
		return rtf.format(diffMins, 'minute');
	}

	const diffHours = Math.round(diffMins / 60);
	if (Math.abs(diffHours) < 24) {
		return rtf.format(diffHours, 'hour');
	}

	const diffDays = Math.round(diffHours / 24);
	if (Math.abs(diffDays) < 7) {
		return rtf.format(diffDays, 'day');
	}

	const diffWeeks = Math.round(diffDays / 7);
	if (Math.abs(diffWeeks) < 4) {
		return rtf.format(diffWeeks, 'week');
	}

	const diffMonths = Math.round(diffDays / 30);
	if (Math.abs(diffMonths) < 12) {
		return rtf.format(diffMonths, 'month');
	}

	const diffYears = Math.round(diffDays / 365);
	return rtf.format(diffYears, 'year');
}
