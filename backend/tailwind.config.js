/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.templ"],
  theme: {
    extend: {
			colors: {
				text: 'var(--text)',
				'middle-text': 'var(--middle-text)',
				'sec-text': 'var(--sec-text)',
				background: 'var(--background)',
				background2: 'var(--background2)',
				'header-bg': 'var(--header-bg)',
				'secondary-action': 'var(--secondary-action)',
				primary: 'var(--primary)',
				secondary: 'var(--secondary)',
				accent: 'var(--accent)',
				'gradient-2': 'var(--gradient-2)',
				'gradient-shadow': 'var(--gradient-shadow)',
				'accent-shadow': 'var(--accent-shadow)',
			},
			fontFamily: {
				'pixel': ['Pixel', 'monospace'],
				'inter': ['Inter', 'sans-serif'],
			}
		},
  },
  plugins: [],
};
