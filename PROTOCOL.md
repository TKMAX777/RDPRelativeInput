# RDPRelativeInput Protocol v3

For the program overview, installation, and usage, see the [README](readme.md).

The client and server must both be updated to v3. The text handshake still uses
`RDPRelativeInput`, `ProtocolVersion:3`, and `Status:OK`; a version mismatch is
rejected before input is processed. The add-in forwards the binary stream.

## Event format

After the handshake, each event has a one-byte opcode followed by little-endian
fields. No padding or newline is included:

| Opcode | Event | Payload | Total bytes |
| --- | --- | --- | --- |
| `80` | Relative mouse move | signed 16-bit X, Y | 5 |
| `81` | Large relative move | signed 32-bit X, Y | 9 |
| `82` | Absolute mouse move | unsigned 16-bit X, Y | 5 |
| `83` / `84` | Key down / up | unsigned 8-bit virtual key | 2 |
| `85` / `86` | Mouse button down / up | unsigned 8-bit button (1, 2, 4) | 2 |
| `87` | Wheel | signed 16-bit delta | 3 |
| `88` | Close | none | 1 |

Opcodes are hexadecimal. Events are sent immediately without a batching timer;
large movement values are preserved without clipping. The receiver handles
fragmented and combined frames, validates values, and terminates the session on
invalid or truncated frames. A failed or short write stops subsequent sends to
avoid corrupting frame boundaries.

## Limits and error handling

Incoming DLL data is held in a synchronized 64 KiB buffer. Oversized input,
invalid handshakes, and disconnects stop the affected channel. Handshake lines
are limited to 256 bytes. Pending native writes retain their buffers until
completion or cancellation.

## Validation and benchmarks

Run `go test ./...` and `go test ./protocol -run '^$' -bench . -benchmem` to check
correctness and compare codec costs. These microbenchmarks do not measure RDP
network latency. Disable `RELATIVE_INPUT_DEBUG` during latency measurements;
debug mode writes input events to disk.
