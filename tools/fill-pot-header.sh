#!/bin/sh
sed \
  -e "s/SOME DESCRIPTIVE TITLE/Translations for Comic Sticks (xkcd-gtk)/" \
  -e "/YEAR THE PACKAGE'S COPYRIGHT HOLDER/d" \
  -e 's#"Report-Msgid-Bugs-To: .*\\n"#"Report-Msgid-Bugs-To: https://github.com/rkoesters/xkcd-gtk/issues\\n"#g' \
  -e '/POT-Creation-Date:/d' \
  -e "s/CHARSET/UTF-8/"
