```markdown
# Design System Document: The Scholarly Forum

## 1. Overview & Creative North Star: "The Digital Curator"
The objective of this design system is to move away from the cluttered, chaotic nature of traditional forums and toward a "Digital Curator" experience. We are treating community discussions as editorial content. The "Creative North Star" is **Architectural Clarity**—the interface should feel like a modern library or a high-end gallery, where the architecture (the UI) is silent, and the art (the community content) is the focus.

To achieve this, we reject the rigid, "boxed-in" layout of 2010-era forums. Instead, we utilize intentional white space, tonal depth layering, and sophisticated typography scales to create a sense of breathing room and intellectual prestige.

---

## 2. Colors & Surface Philosophy
The palette uses deep, authoritative navy tones (`primary`) contrasted against soft, lavender-tinted neutrals (`surface`). 

### The "No-Line" Rule
**Strict Guideline:** 1px solid borders are prohibited for sectioning or containment. 
Boundaries are defined exclusively through background shifts. To separate the sidebar from the feed, place a `surface-container-low` panel against a `surface` background. If you need more definition, use `surface-container-high`.

### Surface Hierarchy & Nesting
Treat the UI as a series of physical layers of fine vellum.
- **Base Layer:** `surface` (#fbf8ff) for the main background.
- **Sectioning:** `surface-container-low` (#f5f2fb) for broad layout regions (e.g., the feed area).
- **Interactive Cards:** `surface-container-lowest` (#ffffff) for the post cards themselves to make them "pop" forward naturally.
- **Navigation/Modals:** `surface-container-highest` (#e4e1ea) for elements that require the highest priority.

### The "Glass & Gradient" Rule
For "Join" or "New Post" actions, we move beyond flat color. 
- **Signature CTAs:** Use a subtle gradient from `primary` (#000666) to `primary_container` (#1a237e) at a 135-degree angle. This adds "soul" and a tactile, premium feel.
- **Floating Navigation:** Use Glassmorphism for the top navigation bar: `surface` at 80% opacity with a `backdrop-blur` of 20px.

---

## 3. Typography: The Editorial Voice
The system pairs **Manrope** (Display/Headlines) for a modern, geometric authority with **Inter** (Body/Labels) for world-class legibility.

- **Display-LG (Manrope, 3.5rem):** Use sparingly for hero community titles or major landing headers.
- **Headline-SM (Manrope, 1.5rem):** The standard for Post Titles. The geometric nature of Manrope gives forum threads a "published" feel.
- **Body-LG (Inter, 1rem):** The primary reading size for forum posts. Set with a generous line-height (1.6) to prevent eye fatigue during long-form discussions.
- **Label-MD (Inter, 0.75rem):** Used for metadata—timestamps, category tags, and user handles. Use `on_surface_variant` (#454652) to ensure a clear hierarchy against the body text.

---

## 4. Elevation & Depth
We eschew traditional drop shadows in favor of **Tonal Layering**.

- **The Layering Principle:** Depth is achieved by "stacking." A `surface-container-lowest` card placed on a `surface-container-low` background creates a soft, natural lift that mimics high-end stationery.
- **Ambient Shadows:** Only for floating elements (like dropdowns or modals). Use a highly diffused shadow: `box-shadow: 0 12px 40px rgba(27, 27, 33, 0.06);`. Note the color is a tint of `on_surface`, not pure black.
- **The Ghost Border Fallback:** If accessibility requirements demand more contrast, use a "Ghost Border": `outline-variant` (#c6c5d4) at 15% opacity. It should be felt, not seen.

---

## 5. Components

### Post Cards
- **Structure:** Use `surface-container-lowest` with a `xl` (0.75rem) corner radius.
- **Spacing:** Minimum 24px internal padding (Spacing Scale).
- **Separation:** Forbid the use of divider lines between posts. Use 16px of vertical white space to separate cards.

### Buttons (The Interaction Set)
- **Primary (Join/Post):** Gradient from `primary` to `primary_container`. Text: `on_primary`. Radius: `full`.
- **Secondary (Filter/Category):** `secondary_container` background with `on_secondary_container` text.
- **Tertiary (Like/Dislike):** Ghost style. No background or border. Use `on_surface_variant` for the icon. On hover, transition to `tertiary_container` (#4f2400) at 10% opacity.

### Navigation Chips
- **Selection Chips:** Used for categories (e.g., "Design," "Engineering"). Use `secondary_fixed` (#cfe6f2) with `on_secondary_fixed` (#071e27) text. Use `md` (0.375rem) corner radius.

### Input Fields
- **Search/Post Area:** Background `surface_container_low`. On focus, transition to `surface_container_lowest` with a 2px `primary` bottom-border only. This mimics a "notepaper" feel.

---

## 6. Do's and Don'ts

### Do
- **Do** prioritize white space. If a layout feels "tight," double the padding before adding a border.
- **Do** use the `tertiary` (#301300) and `tertiary_container` (#4f2400) accents for "Hot" or "Trending" labels to provide a warm contrast to the cool blues.
- **Do** use `manrope` for numbers (upvote counts). Its geometric clarity is superior for data.

### Don't
- **Don't** use 100% black text. Always use `on_surface` (#1b1b21) for body copy to maintain a premium, ink-on-paper look.
- **Don't** use "Card Shadows" on every element. If everything floats, nothing is important. Use background color shifts as the primary tool.
- **Don't** use harsh, high-contrast dividers. If you must separate metadata, use a 4px bullet point (•) in `outline_variant`.

---

## 7. Signature Elements
- **The "Activity Glow":** For active users or new notifications, use a small 8px circle with a subtle outer glow using `surface_tint` (#4c56af).
- **Asymmetric Accents:** For the main community landing page, consider an asymmetric hero section where the `display-lg` text is offset to the left, and a `surface_container_high` decorative element overlaps the background, breaking the standard container grid.```