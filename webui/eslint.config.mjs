/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * @author ENDERZOMBI102 <enderzombi102.end@gmail.com> 2024
 * @description Quick and dirty `eslint` config to better conform to the Prof's requests and style.
 */
import vue from 'eslint-plugin-vue';

// noinspection JSUnusedGlobalSymbols
export default [
	// Ignori globali: file di terze parti e artefatti generati, non codice del
	// progetto. Senza questi, `eslint .` analizza la copia vendored di Bootstrap
	// (public/bootstrap/js) e il binario di Yarn (.yarn/releases), segnalando
	// problemi che non appartengono al sorgente dell'applicazione.
	// Le regole sotto restano identiche alla configurazione fornita dal progetto.
	{
		ignores: [
			'.yarn/**',
			'dist/**',
			'public/bootstrap/**'
		]
	},
	... vue.configs[ "flat/recommended" ],
	{
		rules: {
			'vue/multi-word-component-names': 'off',
			'vue/max-attributes-per-line': 'off',
			'vue/require-default-prop': 'off',
			'vue/singleline-html-element-content-newline': 'off'
		}
	},
];
