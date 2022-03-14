#!/bin/sh
sed \
  -e "s/SOME DESCRIPTIVE TITLE/Translations for Comic Sticks (xkcd-gtk)/" \
  -e "s/YEAR THE PACKAGE'S COPYRIGHT HOLDER/2015-2022 Ryan Koesters/" \
  -e '/^"POT-Creation-Date: .*\\n"$/d' \
  -e "s/CHARSET/UTF-8/"
