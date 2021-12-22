# **Библиотека генерации игр CoinCup**

## **Crash**

### **Проверка игры**

Для проверки честности текущей игры, достаточно запомнить "Серийный номер", отображаемый над графиком.

После окончания игры, кликните по коэффициенту прошедшей игры в ленте над графиком. В открывшемся окне будет указан "Серийный номер", "Дата создания" и ссылка на random.org, где можно сверить данные.

Пример игры №7734

![информация о игре](/assets/crash_form.png)

Далее по ссылке на random.org можно увидеть дату генерации и серийный номер.

![генерация на random.org](/assets/crash_random_org.png)

**Внимание!** Дата и время создания игры указано в часовом поясе вашего устройства.

### **Генерация результата**

Как только появляется отсчет времени до начала игры, происходит генерация результата игры.

Число генерируется на сайте [random.org](https://random.org) в виде десятичной дроби в диапазоне от 0 до 1 с тремя знаками после запятой.

Далее полученное число конвертируется в коэффициент краша по следующей формуле.

```
func crashFloor(value float64) float64 {
	result := 0.05 + 0.95/(1-value)
	if int(math.Floor(result*100))%33 == 0 {
		result = 1
	} else {
		result = math.Round(result*100) / 100
	}
	return result
}
```

По итогу мы получаем коэффициент краша result от 1х до 999х.

## **Double**

### **Проверка игры**

Для проверки честности текущей игры, достаточно запомнить "Серийный номер", отображаемый над колесом.

После окончания игры, кликните по кнопке "Проверка игры", в открывшемся окне будет указан "Серийный номер", "Дата создания" и ссылка на random.org, где можно сверить данные.

Сгенерированное число соответствует следующему коэффициенту (цвету):

- х2 - 1, 3, 5, 7, 9, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 53
- х3 - 0, 4, 6, 8, 10, 15, 17, 19, 25, 27, 29, 31, 33, 39, 41, 43, 49
- х5 - 2, 11, 13, 21, 23, 35, 37, 45, 47, 51
- х50 - 52

### **Генерация результата**

Как только появляется отсчет времени до начала игры, происходит генерация результата игры.

На сайте [random.org](https://random.org) генерируется целое число от 0 до 53.

Полученное число является индексом в ниже следующем массиве множителей.

```
var doubleCoefficients = []uint8{
	3, 2, 5, 2, 3, 2, 3, 2, 3, 2, 3, 5, 2, 5, 2, 3, 2, 3, 2, 3, 2, 5, 2, 5, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 5, 2, 5, 2, 3, 2, 3, 2, 3, 2, 5, 2, 5, 2, 3, 2, 5, 50, 2,
}
```

Данный массив является копией цветов из картинки колеса.

![колесо](/assets/wheel.svg)

## **Mines**

### **Почему генерация не на random.org?**

- Random.org имеет ограничение на количество генераций, а именно 10 запросов в секунду. Для комфортной игры множества человек, такое решение не подходит.

- Так же для генерации каждого результата игры необходима оплата, что требует повышения минимальной ставки на игру.

### **Проверка игры**

В любой момент времени над полем справа доступка кнопка "Проверка игры", при клике на которую открывается окно с "Хэшем игры" и "Датой создания игры".

Результат генерируется, как только вы впервые открыли игру. В дальнейшем генерация происходит после завершения текущей игры.

В хэше зашифрован результат игры, методом [SHA512](https://emn178.github.io/online-tools/sha512.html).

По завершению игры в окне проверки будет доступен "Результат игры".

Для проверки того, что хэш игры не был подменен, вы можете вставить "Результат игры" на любом сайте с SHA512 и сверить полученный хэш с указанным ранее в поле "Хэш игры".

Результат игры представлен в виде строки разделенной символом **|**.

Например:

> 2iouaxvcuetdt2hj291s7fg1pgu5zk68xyk5|**3|22|23|17|19|8|12|5|11|1|24|10|20|25|6|4|9|13|14|15|7|18|2|16|21**|6wmifxp6i4cil4bi0tju1m2hsmy6h71jg37n

По левой и правой части размещена соль, случайная последовательность цифр и букв, обеспечивающая безопасность результата от перебора.

В центре строки размещены номера полей, где размещены мины.

Допустим мы играли в игру с 5 минами, тогда выбираются 5 первых чисел, как места их расположения.

> 3 22 23 17 19

Нумерация ячеек идет от 1 до 25 по принципу слева направо сверху вниз.

### **Генерация результата**

Генерация расположения мин работает по следующему алгоритму

```
base := []uint8{
  1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
  11, 12, 13, 14, 15, 16, 17, 18,
  19, 20, 21, 22, 23, 24, 25,
}
places := make([]uint8, 25)
for i := 0; i < 25; i++ {
  baseLength := len(base)
  r := rand.Intn(baseLength)
  places[i] = base[r]
  base[r] = base[baseLength-1]
  base = base[:baseLength-1]
}
```

Далее генерируем соль и соединяем все вместе

```
leftSeed := make([]rune, lettersLength)
rightSeed := make([]rune, lettersLength)
for i := range leftSeed {
  leftSeed[i] = letters[rand.Intn(lettersLength)]
  rightSeed[i] = letters[rand.Intn(lettersLength)]
}

join := joinUint8(places, "|")

result := fmt.Sprintf("%s|%s|%s", string(leftSeed), join, string(rightSeed))
```

По итогу получаем result строку и хэшируем её в SHA512

## **Dice**

### **Почему генерация не на random.org?**

- Random.org имеет ограничение на количество генераций, а именно 10 запросов в секунду. Для комфортной игры множества человек, такое решение не подходит.

- Так же для генерации каждого результата игры необходима оплата, что требует повышения минимальной ставки на игру.

### **Проверка игры**

В любой момент времени над полем справа доступка кнопка "Проверка игры", при клике на которую открывается окно с "Хэшем игры" и "Датой создания игры".

Результат генерируется, как только вы впервые открыли игру. В дальнейшем генерация происходит после завершения текущей игры.

В хэше зашифрован результат игры, методом [SHA512](https://emn178.github.io/online-tools/sha512.html).

По завершению игры в истории игр будет доступена прошедшая игра с кнопкой проверки,
в окне проверки игры появится поле "Результат игры".

Результат игры представлен в виде строки разделенной символом **|**.

Например:

> pkr4pvp7r1wem5nzjw0uc5zqxyyzc9dl868u|**476279**|fwsmm635t914z5i6ut394tbwssyhf26whqzc

По левой и правой части размещена соль, случайная последовательность цифр и букв, обеспечивающая безопасность результата от перебора.

В центре строки размещено выигрышное число.

### **Генерация результата**

Генерация выигрышного числа работает по следующему алгоритму

```
leftSeed := make([]rune, lettersLength)
rightSeed := make([]rune, lettersLength)
for i := range leftSeed {
  leftSeed[i] = letters[rand.Intn(lettersLength)]
  rightSeed[i] = letters[rand.Intn(lettersLength)]
}

value := l.rand.Intn(1000000)
result := fmt.Sprintf("%s|%d|%s", string(leftSeed), value, string(rightSeed))
hash := sha512.New()
hash.Write([]byte(result))
resultHash := fmt.Sprintf("%x", hash.Sum(nil))
```

По итогу мы получаем result, как исходную строку содержащую выигрышное число и resultHash, захешированный результат игры, который доступен до начала игры.
