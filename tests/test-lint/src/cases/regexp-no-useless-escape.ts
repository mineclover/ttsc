/**
 * Verifies regexp/no-useless-escape: backslash before a non-special character.
 *
 * Pins the alias branch that flags useless escapes inside regex literals, such
 * as `/\a/` where the leading backslash carries no semantic meaning and is
 * usually a typo for an intended escape sequence.
 *
 * 1. Declare a regex literal with a `\a` escape that is not a special escape.
 * 2. Assert it is flagged.
 */
// expect: regexp/no-useless-escape error
const escape = /\a/;

JSON.stringify(escape);
