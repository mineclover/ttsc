/**
 * Verifies jsx-a11y/media-has-caption: `<audio>`/`<video>` need
 * captions.
 *
 * Pins the "no caption track" branch — without a `<track
 * kind="captions">` child, deaf and hard-of-hearing users have no
 * access to the audio content, so the rule rejects the media element.
 *
 * 1. Render a `<video>` with no `<track kind="captions">` child.
 * 2. Lint flags the missing caption track.
 */
export const X = () => <video src="/clip.mp4" />;
