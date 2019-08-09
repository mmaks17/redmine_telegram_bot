package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	_ "github.com/lib/pq"
)

type Issues struct {
	XMLName    xml.Name `xml:"issues"`
	Text       string   `xml:",chardata"`
	TotalCount string   `xml:"total_count,attr"`
	Offset     string   `xml:"offset,attr"`
	Limit      string   `xml:"limit,attr"`
	Type       string   `xml:"type,attr"`
	Issue      []struct {
		Text    string `xml:",chardata"`
		ID      string `xml:"id"`
		Project struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"project"`
		Tracker struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"tracker"`
		Status struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"status"`
		Priority struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"priority"`
		Author struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"author"`
		AssignedTo struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"assigned_to"`
		Subject        string `xml:"subject"`
		Description    string `xml:"description"`
		StartDate      string `xml:"start_date"`
		DueDate        string `xml:"due_date"`
		DoneRatio      string `xml:"done_ratio"`
		EstimatedHours string `xml:"estimated_hours"`
		CustomFields   struct {
			Text        string `xml:",chardata"`
			Type        string `xml:"type,attr"`
			CustomField []struct {
				Text     string `xml:",chardata"`
				ID       string `xml:"id,attr"`
				Name     string `xml:"name,attr"`
				Multiple string `xml:"multiple,attr"`
				Value    struct {
					Text string `xml:",chardata"`
					Type string `xml:"type,attr"`
				} `xml:"value"`
			} `xml:"custom_field"`
		} `xml:"custom_fields"`
		CreatedOn string `xml:"created_on"`
		UpdatedOn string `xml:"updated_on"`
		ClosedOn  string `xml:"closed_on"`
		Parent    struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
		} `xml:"parent"`
	} `xml:"issue"`
}

var (
	firtuse bool
	// глобальная переменная в которой храним токен
	telegramBotToken string
	L                Issues
	url              string
	usernames        string
	version          string
	apikey           string
	db               *sql.DB
	err              error
	postgres_uri     string
)

type BotanMessage struct {
	Text   string
	ChatId int
}

func init() {
	telegramBotToken = os.Getenv("TG_TOKEN")
	if telegramBotToken == "" {
		log.Print("-telegrambottoken is required")
		os.Exit(1)
	}

	url = os.Getenv("REDMINE_URL")
	apikey = os.Getenv("REDMINE_API")
	postgres_uri = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("PG_HOST"), os.Getenv("PG_PORT"), os.Getenv("PG_USER"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_DB"))

	log.Println(postgres_uri)
	firtuse = true
	version = "0.2.2"

}
func getIssues(userid string) Issues {
	req, err := http.NewRequest("GET", url+"/issues.xml?offset=0&limit=1000&assigned_to_id="+userid, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("X-Redmine-Api-Key", apikey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//responseString := string(responseData)

	var ctask Issues
	jsonErr := xml.Unmarshal(responseData, &ctask)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return ctask
}

//переписать нахуй, уже не так
func issueList(t Issues) string {
	var ts string
	//log.Println(len(ctask.Issue))

	for _, x := range t.Issue {

		//	log.Println(x.Subject + "\n\r")

		ts += x.Status.Name + ": " + "[" + x.Subject + "](" + url + "/issues/" + x.ID + ")" + "\n\r"

	}

	return ts
}
func getIDS() []int {
	a := make([]int, 1)
	rows, err := db.Query("SELECT distinct(id) FROM users")
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		a = append(a, id)
	}
	return a

}
func getChatid(id int) int64 {
	if id == 0 {
		return 0
	}
	row := db.QueryRow("SELECT chat FROM users where id=$1 limit 1", id)
	log.Println(" в функции  номер чата для " + strconv.Itoa(id))

	var rez int64
	err = row.Scan(&rez)
	if err != nil {
		panic(err)
	}
	fmt.Println(rez)
	return rez
}
func loopUpdate() {
	var tablename string
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	for {
		if firtuse {
			tablename = "tasks"

			_, err := db.Exec("delete from tasks")
			if err != nil {
				log.Println("Не удалось  почистить  при запуске tasks  " + err.Error())
			}
			_, err2 := db.Exec("delete from tasks2")
			if err2 != nil {
				log.Println("Не удалось  почистить  при запуске tasks2  " + err2.Error())
			}

		} else {
			tablename = "tasks2"
		}

		// получили все уникальные ид пользователей в базе
		s := getIDS()

		fmt.Printf("%+v\n", s)
		for _, idi := range s {
			if idi == 0 {
				continue
			}
			//log.Printf("%d\n", idi)

			//проходимся по всем ишью для каждого пользователя  и записываем в базу
			I := getIssues(strconv.Itoa(idi))
			for _, x := range I.Issue {

				//add feature1
				if x.Author.ID == x.AssignedTo.ID {
					continue
				}

				db.Exec("insert into "+tablename+"  values ($1, $2, $3 , $4)", idi, x.ID, x.Status.Name, x.Subject)
				result, err := db.Exec("insert into "+tablename+"  values ($1, $2, $3 , $4)", idi, x.ID, x.Status.Name, x.Subject)
				if err != nil {
					log.Println(err)
				} else {
					fmt.Println(result.RowsAffected()) // количество добавленных строк
				}

			}

			rows, err := db.Query("select distinct(taskid),status,subject from tasks2 where taskid  not in (select taskid from tasks where userid=$1) and userid=$1", idi)
			if err != nil {
				//log.Println(err)
			}
			replstr := ""

			for rows.Next() {
				var taskid string
				var status string
				var subject string
				err = rows.Scan(&taskid, &status, &subject)
				replstr += "Задача :" + status + " \n\r" + subject + "  \n\r  " + url + "/issues/" + taskid + "\n\r"
				log.Println("Задача :" + status + " " + subject + "   " + url + "/issues/" + taskid)

			}
			if firtuse {
				firtuse = false
				continue
			}

			//  удалить из  tasks  задания для текущего пользователя.
			result, err := db.Exec("delete from tasks where userid = $1", idi)
			if err != nil {
				log.Println("Не удалось удалить данные из tasks для " + strconv.Itoa(idi))
			}
			fmt.Println(result.RowsAffected()) // количество удаленных строк

			//перегоняем новые данные из tasks2    в  tasks
			result2, err2 := db.Exec("insert into tasks select * from tasks2 where userid= $1", idi)
			if err2 != nil {
				log.Println("Не удалось перенести данные в  tasks" + strconv.Itoa(idi))
			}
			fmt.Println(result2.RowsAffected()) // количество удаленных строк

			//удаляем из  tasks2  уже устаревшие заиси
			result3, err3 := db.Exec("delete from tasks2 where userid = $1", idi)
			if err3 != nil {
				log.Println("Не удалось удалить данные из tasks2 для " + strconv.Itoa(idi))
			}
			fmt.Println(result3.RowsAffected()) // количество удаленных строк

			log.Println(" получаем номер чата для " + strconv.Itoa(idi))
			msg := tgbotapi.NewMessage(getChatid(idi), " ")
			msg.ParseMode = "markdown"

			msg.Text += replstr
			log.Println(msg.Text)
			bot.Send(msg)

		}

		time.Sleep(15 * time.Second)
	}

}

func getmyid(user string) string {
	row := db.QueryRow("select id from users where name = $1", user)
	var rez string
	err = row.Scan(&rez)
	if err != nil {
		panic(err)
	}
	fmt.Println(rez)
	return rez
}

func sendSpam(message string) {
	// получили все уникальные ид пользователей в базе
	bot, _ := tgbotapi.NewBotAPI(telegramBotToken)
	s := getIDS()

	fmt.Printf("%+v\n", s)
	for _, ids := range s {
		if ids == 0 {
			continue
		}
		getChatid(ids)

		msg := tgbotapi.NewMessage(getChatid(ids), " ")
		msg.ParseMode = "markdown"

		msg.Text = message
		log.Println(msg.Text)
		bot.Send(msg)

	}

}
func main() {

	log.Println("Waiting wgen db up ")
	time.Sleep(10 * time.Second)
	// Connecting to database
	db, err = sql.Open("postgres", postgres_uri)
	if err != nil {
		log.Fatalf("Can't connect to postgresql: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Can't ping database: %v", err)
	}

	// используя токен создаем новый инстанс бота
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	go loopUpdate()

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		// универсальный ответ на любое сообщение
		reply := ""
		if update.CallbackQuery != nil {
			//message := BotanMessage{Text: update.CallbackQuery.Data, ChatId: update.CallbackQuery.Message.From.ID}
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ok, I remember^- "+
				update.CallbackQuery.Data+update.CallbackQuery.ID+string(update.CallbackQuery.From.UserName)))
		}

		if update.Message == nil {
			continue
		}
		// комманда - сообщение, начинающееся с "/"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		keyboard := tgbotapi.InlineKeyboardMarkup{}

		switch update.Message.Command() {

		case "start":
			reply = "  работая с этим ботом вы соглашаетесь с политикой конфидециальности и правилами работы с сервисом. \n\r" +
				"Для  отслеживания задач  введите /set id ,   где  id -   номер вашего пользователя в редмайне. \n\r" +
				"Введите /delete  для того что бы отписаться от подписки на сообщения. \n\r" +
				"Соответствие ид пользователя в редмайне,  и пользователя который его себе установил никак не отслеживается, будье  внимательны и благоразумны. "

		case "delete":
			log.Println(update.Message.Chat.UserName)
			//если текущий id  редмайна уже ктото прослушивает, перезаписываем себе
			//	result, err := db.Exec("insert into users (id, name, chat) values ($1, $2, $3)", update.Message.CommandArguments(), update.Message.Chat.UserName, update.Message.Chat.ID)
			result, err := db.Exec("delete from users where name=$1 ", update.Message.Chat.UserName)

			if err != nil {
				log.Printf("%v\n\r", err)
			} else {
				reply = "  вы отписались "
				fmt.Println(result.RowsAffected()) // количество добавленных строк
			}

		case "set":
			s1 := update.Message.CommandArguments()
			if s1 != "" {
				str, _ := strconv.ParseInt(s1, 10, 64)
				if str <= 0 {
					reply = " не для групповых чатов, пишите в личку"
				} else {
					log.Println(update.Message.Chat.UserName)
					//если текущий id  редмайна уже ктото прослушивает, перезаписываем себе
					//	result, err := db.Exec("insert into users (id, name, chat) values ($1, $2, $3)", update.Message.CommandArguments(), update.Message.Chat.UserName, update.Message.Chat.ID)
					result, err := db.Exec("insert into users (id, name, chat) values ($1, $2, $3)  ON CONFLICT (id) DO  UPDATE   SET  name=$2, chat=$3 ", update.Message.CommandArguments(), update.Message.Chat.UserName, update.Message.Chat.ID)

					if err != nil {
						panic(err)
					} else {
						reply = " успешно добавлено "
						fmt.Println(result.RowsAffected()) // количество добавленных строк
					}

				}

			} else {
				reply = " не для групповых чатов, пишите в личку"
			}

		case "u":

			rows, err := db.Query("SELECT * FROM users")
			if err != nil {
				log.Println(err)
			}
			s := ""

			for rows.Next() {
				var id int
				var name string
				var chat int
				err = rows.Scan(&id, &name, &chat)
				s += strconv.Itoa(id) + " " + name + " " + strconv.Itoa(chat) + "\n\r"

				fmt.Println("uid | name | chat  ")
				fmt.Printf("%3v | %8v | %6v \n", id, name, chat)
			}
			reply = s

		case "/btn":
			var row []tgbotapi.InlineKeyboardButton
			btn := tgbotapi.NewInlineKeyboardButtonData("text", "data")
			row = append(row, btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			msg.ReplyMarkup = keyboard
			reply = "you are ready ?"

		case "version":
			reply = version

		case "l":
			reply = issueList(getIssues(getmyid(update.Message.Chat.UserName)))
		case "reklama":
			sendSpam(update.Message.CommandArguments())

		default:
			continue
		}

		//log.Printf("%s %s", update.Message.Chat.ID, reply)

		// создаем ответное сообщение
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		// отправляем
		bot.Send(msg)

	}
}
