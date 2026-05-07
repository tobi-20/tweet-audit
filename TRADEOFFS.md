## Architecture Overview

```
Data → Parser → Brain (ProcessTweets) → External IO (AI + Writer)
```
I structured the system  as a pipeline where ProcessTweets is the coordinator, and external dependencies like the AI client and writer are abstracted with interfaces. This keeps business logic separate from external I/O and makes the system easier to test by  being able to swap implementations.


## Concurrency Strategy

**Decision: Sequential execution.**

| Option | Trade-off |
|--------|-----------|
| Sequential | Predictable, checkpoint-safe, rate-limit-friendly |
| Goroutines | Faster, but requires careful batching and complicates recovery |

Goroutines are only justified if you stay safely with strict concurrency limits. For now, sequential wins because external APIs are unpredictable and recoverability beats speed.


## Error Handling

**Approach: Fail-fast + checkpoint.**

- API error → stop execution immediately
- Resume from last saved progress on restart
- Progress is saved **only after a successful response** (never on failure)

```
item 4 fails → saveProgress(3) → restart later → resumes at 4
```

Saving on failure creates ambiguity. Saving only on success guarantees idempotent restarts.

---

## Performance vs. Safety Trade-offs

| Factor | Decision | Reason |
|--------|----------|--------|
| Speed | Low priority | External API can have unpredictable latency due to network speed and other factors |
| API cost | Controlled | Sequential processing helps prevent unnecessary expenses |
| Reliability | High priority | Checkpoint recovery is the core reason |
| Concurrency | None | Adds unnecessary complexity right now |


---

## Why I used Go for this Problem

- **Familiarity** — I am most comfortable and familiar with Go for backend related projects.

- **Simple error model** — `(value, error)` return pairs make failing fast easier and error handling easier.

- **Easy file interactions** — There are various built in functions that make I/O operations easy manipulable

---

## Resumability Design

The checkpoint file holds the last successfully processed index.

```
On start  → read checkpoint → start = checkpoint + 1
On success → write checkpoint → i
On failure → return error → checkpoint unchanged
```

