package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type VkGroup struct {
	Name string
}

type Db struct {
	DBConnection *sql.DB
}

func NewDbConnection() *sql.DB {
	// Получаем путь к БД из переменной окружения или используем значение по умолчанию
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "../glavredus.db"
	}

	// Если путь относительный, преобразуем в абсолютный
	path, err := filepath.Abs(dbPath)
	if err != nil {
		fmt.Println("Error resolving DB path:", err)
		path = dbPath
	}

	log.Printf("Using database path: %s", path)

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}

	return db
}

func LoadSchema(db *sql.DB) {
	deleteVkGroupTable(db)
	createVkgroupSQL := `
		CREATE TABLE IF NOT EXISTS vkgroup (
			name string PRIMARY KEY NOT NULL
		);`

	_, err := db.Exec(createVkgroupSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteVkGroupTable(db *sql.DB) {
	sql := `DROP TABLE IF EXISTS vkgroup;`

	_, err := db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadSchoolVkGroups(db *sql.DB) {
	vkgroups := []VkGroup{
		{"krsel"},
		{"club194809745"},
		{"club214119048"},
		{"club202724280"},
		{"club185982638"},
		{"club205401563"},
		{"club205402681"},
		{"detskisad15"},
		{"16detskiysad"},
		{"club205401551"},
		{"doy19"},
		{"club109060055"},
		{"club205401929"},
		{"club205400972"},
		{"club182072023"},
		{"club195576991"},
		{"club147892228"},
		{"club187951249"},
		{"sadik31krs"},
		{"club205420428"},
		{"gdboy35"},
		{"club205443755"},
		{"dc39spb"},
		{"club170186955"},
		{"gbdou41krspb"},
		{"club216246675"},
		{"club205406349"},
		{"club203026295"},
		{"dc5krs"},
		{"ds51krs"},
		{"gbdouds52"},
		{"club13309436"},
		{"club192983329"},
		{"club205417092"},
		{"club214317110"},
		{"gbdou6kr"},
		{"club42266729"},
		{"club76873688"},
		{"club202836702"},
		{"club202821332"},
		{"ds_65_krs_spb"},
		{"club205400739"},
		{"dou69krasnosel"},
		{"club216939970"},
		{"club205428969"},
		{"club205401911"},
		{"detskiy_sad74"},
		{"ds75spb"},
		{"club202011664"},
		{"club205406444"},
		{"ds78spb"},
		{"club129697643"},
		{"ds80krs"},
		{"club195029092"},
		{"club203610472"},
		{"gbdou83"},
		{"club203812364"},
		{"club205421015"},
		{"club202723926"},
		{"club215846431"},
		// {"istokdetsad"},
		{"club194904593"},
		{"dc9spb"},
		{"children322029"},
		{"dou91krasnosel"},
		{"club205413257"},
		{"gbdou93krasnosel"},
		{"gbdou94"},
		{"gbdou95"},
		{"club227261708"},
		{"club183141138"},
		{"club205420830"},
		{"club193884037"},
		{"club214016041"},
		{"club200294876"},
		{"68rostok"},
		{"club205440005"},
		{"dc50krs_spb"},
		{"club180362982"},
		{"school509spb"},
		{"schoolspb54"},
		{"gym271"},
		{"gim293spb"},
		{"spb.school399"},
		{"club117133342"},
		{"public220312271"},
		{"licey_369"},
		{"licei395"},
		{"public__590"},
		{"club23933409"},
		{"club215520444"},
		{"school200spb"},
		{"rr_school208"},
		{"vr_odod_237"},
		{"sovet247"},
		{"school252spb"},
		{"spbschool262"},
		{"schooll270"},
		{"sch276spb"},
		{"school285spb"},
		{"g2343"},
		{"gbou291"},
		{"school352veteranov151"},
		{"school382spb"},
		{"school383"},
		{"club214266378"},
		{"spbschool390"},
		{"spbgboy391"},
		{"school394spb"},
		{"school414"},
		{"newschool546"},
		{"school547"},
		{"sc548"},
		{"spbschool549"},
		{"co_167"},
		{"school_289"},
		{"gboyshkola131"},
		{"school203spb"},
		{"school217spb"},
		{"newschool219"},
		{"school242spb"},
		{"forestschool275"},
		{"school375spb"},
		{"school380spb"},
		{"sch398"},
		{"club153653196"},
		{"school7spb"},
		{"csridi_geroev"},
		{"ddtks"},
		{"cbzh_cgpv"},
		{"imc_krsel"},
		{"shkrsl"},
		{"guzhakra"},
		{"club88310495"},
		{"club164410468"},
		{"kdk_krasnoselsky"},
		{"krasnoselskiy_kcson"},
		{"pmcligovo"},
		{"cfksiz"},
		{"club171353821"},
		{"gp106"},
		{"club200827129"},
		{"club147440843"},
		{"club215863666"},
		{"club215786855"},
		{"cbs_krlib"},
		{"krocpmsskr"},
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO vkgroup (name) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, vkgroup := range vkgroups {
		if _, err := stmt.Exec(vkgroup.Name); err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
