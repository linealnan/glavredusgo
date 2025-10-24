package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	index "github.com/linealnan/glavredusgo/fts/internal/index"
	vkapi "github.com/romanraspopov/golang-vk-api"
)

type MockGroup struct {
	Name string
}

// UserCity содержит id и название населенного пункта пользователя ВК
// Информация о городе, указанном на странице пользователя в разделе «Контакты».
// Возвращаются следующие поля:
// id (integer) — идентификатор города, который можно использовать для получения его названия с помощью метода database.getCitiesById;
// title (string) — название города.
type UserCity struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// Full-Text Search (FTS)
// Raw Text -> tokenizer->filters->tokens
// https://habr.com/ru/articles/519024/
// https://github.com/akrylysov/simplefts
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("VK_API_TOKEN")
	client, err := vkapi.NewVKClientWithToken(token, nil, true)
	if err != nil {
		log.Fatal(err)
	}
	loadGroupsData(client)
}

func loadGroupsData(client *vkapi.VKClient) {
	var documents []index.Document
	// groups := getGroups()
	groups := getSchoolGroups()
	// groups := getFromKidsGardenGroups()
	log.Printf("Загрузка данных групп\n")
	for _, group := range groups {
		document := getAndIndexedWallPostByGroupName(client, group.Name)
		documents = append(documents, document...)
	}

	// query := "губернатор"
	// query := "Звонок из деканата"
	query := "поручение"
	// query := "эйфория"

	start := time.Now()
	idx := make(index.Index)
	idx.Add(documents)
	log.Printf("Indexed %d documents in %v", len(documents), time.Since(start))

	start = time.Now()
	matchedIDs := idx.Search(query)
	log.Printf("Search found %d documents in %v", len(matchedIDs), time.Since(start))

	for _, id := range matchedIDs {
		// log.Printf("%d\t", id)
		for _, doc := range documents {
			if doc.ID == id {
				log.Printf("%d\t%s\n", id, doc.URL)
			}
		}
	}
}

func getAndIndexedWallPostByGroupName(client *vkapi.VKClient, groupName string) []index.Document {
	var documents []index.Document
	var document index.Document
	var groupsName []string

	groupsName = append(groupsName, groupName)
	group, err := client.GroupsGetByID(groupsName)
	if err != nil {
		log.Fatal(err)
	}

	wall, err := client.WallGet(groupName, 100, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range wall.Posts {
		// log.Printf("Wall post: %v\n", post.Text)
		i, err := strconv.Atoi(strconv.Itoa(group[0].ID) + strconv.Itoa(post.ID))
		if err != nil {
			log.Fatal("Error converting string to int:", err)
		}
		document.ID = i
		document.Text = post.Text

		document.URL = "https://vk.com/" + groupName + "?w=wall-" + strconv.Itoa(group[0].ID) + "_" + strconv.Itoa(post.ID)

		documents = append(documents, document)
	}

	// for _, doc := range documents {
	// 	log.Printf("%d\t%s\n", doc.ID, doc.URL)
	// }
	return documents
}

// https://vk.com/csridi_geroev
// Администрация Красносельского района Санкт-Петербурга
// ГБУ ДО ДДТ КРАСНОСЕЛЬСКОГО РАЙОНА САНКТ-ПЕТЕРБУРГА
// https://vk.com/ddtks
// https://vk.com/cbzh_cgpv
// Администрация Красносельского района Санкт-Петербурга
// ГБУ ИМЦ КРАСНОСЕЛЬСКОГО РАЙОНА САНКТ-ПЕТЕРБУРГА
// https://vk.com/imc_krsel
// Администрация Красносельского района Санкт-Петербурга
// ГБУ СШ КРАСНОСЕЛЬСКОГО РАЙОНА САНКТ ПЕТЕРБУРГА
// https://vk.com/shkrsl
// Администрация Красносельского района Санкт-Петербурга
// КРАСНОСЕЛЬСКОЕ РЖА https://vk.com/guzhakra
// Администрация Красносельского района Санкт-Петербурга
// ОАМ ЦСРИДИ Красносельского района
// https://vk.com/club88310495
// Администрация Красносельского района Санкт-Петербурга
// Отделение ЦСРИДИ г. Красное Cело
// https://vk.com/club164410468
// Администрация Красносельского района Санкт-Петербурга
// СПБ ГБУ "КДК "КРАСНОСЕЛЬСКИЙ"
// https://vk.com/kdk_krasnoselsky
// Администрация Красносельского района Санкт-Петербурга
// СПБ ГБУ "КЦСОН КРАСНОСЕЛЬСКОГО РАЙОНА"
// https://vk.com/krasnoselskiy_kcson
// Администрация Красносельского района Санкт-Петербурга
// СПБ ГБУ "ПМЦ "ЛИГОВО"https://vk.com/pmcligovo
// Администрация Красносельского района Санкт-Петербурга
// СПБ ГБУ "ЦФКС И З КРАСНОСЕЛЬСКОГО РАЙОНА"
// https://vk.com/cfksiz
// Администрация Красносельского района Санкт-Петербурга
// СПБ ГБУ ДО "ДШИ" КРАСНОСЕЛЬСКОГО РАЙОНА
// https://vk.com/club171353821
// Администрация Красносельского района Санкт-Петербурга
// СПБ ГБУЗ "ГОРОДСКАЯ ПОЛИКЛИНИКА №106"
// https://vk.com/gp106
// Администрация Красносельского района Санкт-Петербурга
// СПБ ГБУЗ "ГОРОДСКАЯ ПОЛИКЛИНИКА №91"
// https://vk.com/club200827129
// Администрация Красносельского района Санкт-Петербурга
// СПБ ГБУЗ "ГОРОДСКАЯ ПОЛИКЛИНИКА №93"
// https://vk.com/club147440843
// Администрация Красносельского района Санкт-Петербурга
// СПБ ГБУЗ "КВД №6"https://vk.com/club215863666
// Администрация Красносельского района Санкт-Петербурга
// СПБ ГБУЗ "СТОМАТОЛОГИЧЕСКАЯ ПОЛИКЛИНИКА №28"
// https://vk.com/club215786855
// Администрация Красносельского района Санкт-Петербурга
// СПБ ГБУК"ЦБС КРАСНОСЕЛЬСКОГО РАЙОНА"
// https://vk.com/cbs_krlib
// Администрация Красносельского района Санкт-Петербурга
// ЦПМСС КРАСНОСЕЛЬСКОГО РАЙОНА
// https://vk.com/krocpmsskr

func getGroups() []MockGroup {
	return []MockGroup{
		{Name: "csridi_geroev"},
		{Name: "ddtks"},
		{Name: "imc_krsel"},
		{Name: "shkrsl"},
		{Name: "guzhakra"},
		{Name: "club88310495"},
		{Name: "kdk_krasnoselsky"},
		{Name: "krasnoselskiy_kcson"},
		{Name: "cfksiz"},
		{Name: "club171353821"},
	}
}

func getSchoolGroups() []MockGroup {
	return []MockGroup{
		{Name: "club194809745"},
		{Name: "club214119048"},
		{Name: "club202724280"},
		{Name: "club185982638"},
		{Name: "club205401563"},
		{Name: "club205402681"},
		{Name: "detskisad15"},
		{Name: "16detskiysad"},
		{Name: "club205401551"},
		{Name: "doy19"},
		{Name: "club109060055"},
		{Name: "club205401929"},
		{Name: "club205400972"},
		{Name: "club182072023"},
		{Name: "club195576991"},
		{Name: "club147892228"},
		{Name: "club187951249"},
		{Name: "sadik31krs"},
		{Name: "club205420428"},
		{Name: "gdboy35"},
		{Name: "club205443755"},
		{Name: "dc39spb"},
		{Name: "club170186955"},
		{Name: "gbdou41krspb"},
		{Name: "club216246675"},
		{Name: "club205406349"},
		{Name: "club203026295"},
		{Name: "dc5krs"},
		{Name: "ds51krs"},
		{Name: "gbdouds52"},
		{Name: "club13309436"},
		{Name: "club192983329"},
		{Name: "club205417092"},
		{Name: "club214317110"},
		{Name: "gbdou6kr"},
		{Name: "club42266729"},
		{Name: "club76873688"},
		{Name: "club202836702"},
		{Name: "club202821332"},
		{Name: "ds_65_krs_spb"},
		{Name: "club205400739"},
		{Name: "dou69krasnosel"},
		{Name: "club216939970"},
		{Name: "club205428969"},
		{Name: "club205401911"},
		{Name: "detskiy_sad74"},
		{Name: "ds75spb"},
		{Name: "club202011664"},
		{Name: "club205406444"},
		{Name: "ds78spb"},
		{Name: "club129697643"},
		// {Name: "ds80krs"},
		// {Name: "club195029092"},
		// {Name: "club203610472"},
		// {Name: "gbdou83"},
		// {Name: "club203812364"},
		// {Name: "club205421015"},
		// {Name: "club202723926"},
		// {Name: "club215846431"},
		// {Name: "istokdetsad"},
		// {Name: "club194904593"},
		// {Name: "dc9spb"},
		// {Name: "children322029"},
		// {Name: "dou91krasnosel"},
		// {Name: "club205413257"},
		// {Name: "gbdou93krasnosel"},
		// {Name: "gbdou94"},
		// {Name: "gbdou95"},
		// https://vk.com/club227261708
		// https://vk.com/club183141138
		// https://vk.com/club205420830
		// https://vk.com/club193884037
		// https://vk.com/club214016041
		// https://vk.com/club200294876
		// https://vk.com/68rostok
		// https://vk.com/club205440005
		// https://vk.com/dc50krs_spb
		// https://vk.com/club180362982

		// https://vk.com/school509spb
		// https://vk.com/schoolspb54
		// https://vk.com/gym271
		// https://vk.com/gim293spb
		// https://vk.com/spb.school399
		// https://vk.com/club117133342
		// https://vk.com/public220312271
		// https://vk.com/licey_369
		// https://vk.com/licei395
		// https://vk.com/public__590
		// https://vk.com/club23933409
		// https://vk.com/club215520444
		// https://vk.com/school200spb
		// https://vk.com/rr_school208
		// https://vk.com/vr_odod_237
		// https://vk.com/sovet247
		// https://vk.com/school252spb
		// https://vk.com/spbschool262
		// https://vk.com/schooll270
		// https://vk.com/sch276spb
		// https://vk.com/school285spb
		// https://vk.com/g2343
		// https://vk.com/gbou291
		// https://vk.com/school352veteranov151
		// https://vk.com/school382spb
		// https://vk.com/school383
		// https://vk.com/club214266378
		// https://vk.com/spbschool390
		// https://vk.com/spbgboy391
		// https://vk.com/school394spb
		// https://vk.com/school414
		// https://vk.com/newschool546
		// https://vk.com/school547
	}
}

// Звонок из «деканата»
// https://vk.com/school252spb
// https://vk.com/spbschool262
// https://vk.com/schooll270
// https://vk.com/sch276spb
// https://vk.com/school285spb
// https://vk.com/g2343
// https://vk.com/gbou291
// https://vk.com/school352veteranov151
// https://vk.com/school382spb
// https://vk.com/school383
// https://vk.com/club214266378
// https://vk.com/spbschool390
// https://vk.com/spbgboy391
// https://vk.com/school394spb
// https://vk.com/school414
// https://vk.com/newschool546
// https://vk.com/school547

// https://vk.com/school509spb
// https://vk.com/schoolspb54
// https://vk.com/gym271
// https://vk.com/gim293spb
// https://vk.com/spb.school399
// https://vk.com/club117133342
// https://vk.com/public220312271
// https://vk.com/licey_369
// https://vk.com/licei395
// https://vk.com/public__590
// https://vk.com/club23933409
// https://vk.com/club215520444
// https://vk.com/school200spb
// https://vk.com/rr_school208
// https://vk.com/vr_odod_237
// https://vk.com/sovet247
func getFromDecanatGroups() []MockGroup {
	return []MockGroup{
		{Name: "school252spb"},
		{Name: "spbschool262"},
		{Name: "schooll270"},
		{Name: "sch276spb"},
		{Name: "school285spb"},
		{Name: "g2343"},
		{Name: "gbou291"},
		{Name: "school352veteranov151"},
		{Name: "school382spb"},
		{Name: "school383"},
		{Name: "club214266378"},
		{Name: "spbschool390"},
		{Name: "spbgboy391"},
		{Name: "school394spb"},
		{Name: "school414"},
		{Name: "newschool546"},
		{Name: "school547"},
		{Name: "school509spb"},
		{Name: "schoolspb54"},
		{Name: "gym271"},
		{Name: "gim293spb"},
		{Name: "spb.school399"},
		{Name: "club117133342"},
		{Name: "public220312271"},
		{Name: "licey_369"},
		{Name: "licei395"},
		{Name: "public__590"},
		{Name: "club23933409"},
		{Name: "club215520444"},
		{Name: "school200spb"},
		{Name: "rr_school208"},
		{Name: "vr_odod_237"},
		{Name: "sovet247"},
	}
}

// Сады
// https://vk.com/club194809745
// https://vk.com/club214119048
// https://vk.com/club202724280
// https://vk.com/club185982638
// https://vk.com/club205401563
// https://vk.com/club205402681
// https://vk.com/detskisad15
// https://vk.com/16detskiysad
// https://vk.com/club205401551
// https://vk.com/doy19
// https://vk.com/club109060055
// https://vk.com/club205401929
// https://vk.com/club205400972
// https://vk.com/club182072023
// https://vk.com/club195576991
// https://vk.com/club147892228
// https://vk.com/club187951249
// https://vk.com/sadik31krs
// https://vk.com/club205420428
// https://vk.com/gdboy35
// https://vk.com/club205443755
// https://vk.com/dc39spb
// https://vk.com/club170186955
// https://vk.com/gbdou41krspb
// https://vk.com/club216246675
// https://vk.com/club205406349
// https://vk.com/club203026295
// https://vk.com/dc5krs
// https://vk.com/ds51krs
// https://vk.com/gbdouds52
// https://vk.com/club13309436
// https://vk.com/club192983329
// https://vk.com/club205417092
// https://vk.com/club214317110
// https://vk.com/gbdou6kr
// https://vk.com/club42266729
// https://vk.com/club76873688
// https://vk.com/club202836702
// https://vk.com/club202821332
// https://vk.com/ds_65_krs_spb
// https://vk.com/club205400739
// https://vk.com/dou69krasnosel
// https://vk.com/club216939970
// https://vk.com/club205428969
// https://vk.com/club205401911
// https://vk.com/detskiy_sad74
// https://vk.com/ds75spb
// https://vk.com/club202011664
// https://vk.com/club205406444
// https://vk.com/ds78spb
// https://vk.com/club129697643
// https://vk.com/ds80krs
// https://vk.com/club195029092
// https://vk.com/club203610472
// https://vk.com/gbdou83
// https://vk.com/club203812364
// https://vk.com/club205421015
// https://vk.com/club202723926
// https://vk.com/club215846431
// https://vk.com/istokdetsad
// https://vk.com/club194904593
// https://vk.com/dc9spb
// https://vk.com/children322029
// https://vk.com/dou91krasnosel
// https://vk.com/club205413257
// https://vk.com/gbdou93krasnosel
// https://vk.com/gbdou94
// https://vk.com/gbdou95
// https://vk.com/club227261708
// https://vk.com/club183141138
// https://vk.com/club205420830
// https://vk.com/club193884037
// https://vk.com/club214016041
// https://vk.com/club200294876
// https://vk.com/68rostok
// https://vk.com/club205440005
// https://vk.com/dc50krs_spb
// https://vk.com/club180362982

func getFromKidsGardenGroups() []MockGroup {
	return []MockGroup{
		{Name: "club194809745"},
		{Name: "club214119048"},
		{Name: "club202724280"},
		{Name: "club185982638"},
		{Name: "club205401563"},
		{Name: "club205402681"},
		{Name: "detskisad15"},
		{Name: "16detskiysad"},
		{Name: "club205401551"},
		{Name: "doy19"},
		{Name: "club109060055"},
		{Name: "club205401929"},
		{Name: "club205400972"},
		{Name: "club182072023"},
		{Name: "club195576991"},
		{Name: "club147892228"},
		{Name: "club187951249"},
		{Name: "sadik31krs"},
		{Name: "club205420428"},
		{Name: "gdboy35"},
		{Name: "club205443755"},
		{Name: "dc39spb"},
		{Name: "club170186955"},
		{Name: "gbdou41krspb"},
		{Name: "club216246675"},
		{Name: "club205406349"},
		{Name: "club203026295"},
		{Name: "dc5krs"},
		{Name: "ds51krs"},
		{Name: "gbdouds52"},
		{Name: "club13309436"},
		{Name: "club192983329"},
		{Name: "club205417092"},
		{Name: "club214317110"},
		{Name: "gbdou6kr"},
		{Name: "club42266729"},
		{Name: "club76873688"},
		{Name: "club202836702"},
	}
}
