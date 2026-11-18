//
// Copyright (c) 2022-2026 Winlin
//
// SPDX-License-Identifier: MIT
//
import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import {Locale} from "./utils";

// See https://react.i18next.com/guides/multiple-translation-files
import resources from "./resources/locale.json";
console.log(`i18n loads resources ${JSON.stringify(Object.keys(resources))}`);

Locale.load();

i18n
  .use(initReactI18next) // Passes i18n down to react-i18next.
  .init({
    resources,
    lng: Locale.current(), // See https://www.i18next.com/overview/configuration-options#languages-namespaces-resources
    // i18n.changeLanguage changes the language later: https://www.i18next.com/overview/api#changelanguage

    interpolation: {
      escapeValue: false // React already escapes values against XSS.
    }
  });

export default i18n;

