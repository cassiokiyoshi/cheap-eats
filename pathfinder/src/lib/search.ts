import { base } from '$app/paths';
import type { SearchRequest, StreamEvent } from '$lib/types';

type EventHandler = (event: StreamEvent) => void;

export async function streamSearch(
	request: SearchRequest,
	signal: AbortSignal,
	onEvent: EventHandler
): Promise<void> {
	const response = await fetch(`${base}/api/search`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(request),
		signal
	});

	if (!response.ok) {
		throw new Error(await readErrorMessage(response));
	}

	if (!response.body) {
		throw new Error('The server returned an empty response.');
	}

	const reader = response.body.getReader();
	const decoder = new TextDecoder();
	let buffer = '';

	while (true) {
		const { value, done } = await reader.read();

		if (done) {
			buffer += decoder.decode();
			break;
		}

		buffer += decoder.decode(value, { stream: true });
		buffer = processCompleteLines(buffer, onEvent);
	}

	processFinalLine(buffer, onEvent);
}

function processCompleteLines(
	buffer: string,
	onEvent: EventHandler
): string {
	const lines = buffer.split('\n');
	const incompleteLine = lines.pop() ?? '';

	for (const line of lines) {
		processLine(line, onEvent);
	}

	return incompleteLine;
}

function processFinalLine(
	buffer: string,
	onEvent: EventHandler
): void {
	if (buffer.trim() !== '') {
		processLine(buffer, onEvent);
	}
}

function processLine(
	line: string,
	onEvent: EventHandler
): void {
	const trimmedLine = line.trim();

	if (trimmedLine === '') {
		return;
	}

	const event = JSON.parse(trimmedLine) as StreamEvent;
	onEvent(event);
}

async function readErrorMessage(response: Response): Promise<string> {
	try {
		const event = (await response.json()) as Partial<StreamEvent>;

		if ('message' in event && typeof event.message === 'string') {
			return event.message;
		}
	} catch {
		// Use the fallback message below.
	}

	return `Search request failed with status ${response.status}.`;
}
