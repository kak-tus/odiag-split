# odiag-split

Разбиение лога OpenDiag на части, полностью загружаемые бесплатной версией DiagView (850 фреймов).

Программа вычитывает все логи из указанной папки. Если лог нужно разбить - он разбивается на новые файлы. Старый файл переименовывается и сохраняется. Если лог уже был разбит или его не нужно разбивать - он не меняется.

Запуск

```
odiag-split <путь к папке с файлами>
```
