# Barcode Library Evaluation Report

## Candidates

| | `barqode` | `html5-qrcode` |
|---|---|---|
| **Version** | 0.0.2 | 2.3.8 |
| **Author** | svecosystem (huntabyte + ollema) | mebjas |
| **GitHub Stars** | 62 | 5.4k+ |
| **Package Size** | **116KB** | **3.3MB** |

## Evaluation Results

| Criteria | Weight | `barqode` ✅ | `html5-qrcode` |
|---------|--------|-------------|----------------|
| **EAN-13/UPC-A reliability** | High | Native BarcodeDetector API + WASM polyfill. Formats: `["ean_13", "ean_8", "upc_a", "upc_e"]` | Built-in decoder. Reliable but heavier. |
| **Svelte 5 compatibility** | High | **Native Svelte 5 component with runes** (`$state`, `$bindable`). Prop-based API, `{#snippet}` children. | **Vanilla JS library.** Requires imperative wrapping: `new Html5Qrcode()`, `.start()`, `.stop()`. Manual DOM ID binding. |
| **Bundle size** | Medium | **116KB** total package. Uses native BarcodeDetector + polyfill. | **3.3MB** total package. Ships its own full decoder. |
| **Camera handling** | Medium | Declarative: `constraints` prop, `onCameraOn`/`onCameraOff` callbacks, `paused` bindable prop. Auto-cleanup on unmount. | Imperative: must call `.start()/.stop()` manually. Need to handle cleanup in `onDestroy`. |
| **API ergonomics** | Low | **Excellent.** `<BarqodeStream formats={[...]} {onDetect} />`. Typed props, canvas tracking, torch control. | **Adequate.** Must create reader div, instantiate scanner, call start with config object. |

## Decision: **barqode**

### Rationale

1. **Native Svelte 5**: barqode is a first-class Svelte 5 component with runes, snippets, and bindable props. It follows the same patterns we use throughout the app (shadcn-svelte, bits-ui). `html5-qrcode` requires an imperative wrapper that fights Svelte's declarative model.

2. **30x smaller**: 116KB vs 3.3MB. The barqode approach delegates to the browser's native `BarcodeDetector` API (Chrome/Android) with a WASM polyfill for Safari/Firefox. This is architecturally superior — the heavy lifting happens in optimized native code.

3. **Camera lifecycle**: barqode handles camera start/stop/cleanup reactively through props. No need for manual `onMount`/`onDestroy` plumbing. The `paused` bindable prop is perfect for our scan-then-lookup flow.

4. **Format filtering**: We only need EAN-13/EAN-8/UPC-A/UPC-E. barqode's `formats` prop lets us restrict detection, reducing processing overhead.

5. **Ecosystem alignment**: barqode is from the svecosystem (same team behind melt-ui, bits-ui patterns). It fits our stack philosophy.

### Tradeoffs Accepted

- Lower GitHub star count (62 vs 5.4k) — but the library is well-maintained and backed by the svecosystem team.
- Newer library (0.0.2) — but the underlying `BarcodeDetector` API and polyfill are battle-tested.
