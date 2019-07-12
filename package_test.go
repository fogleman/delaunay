package delaunay

import (
	"math"
	"math/rand"
	"testing"
)

func uniform(n int, rnd *rand.Rand) []Point {
	points := make([]Point, n)
	for i := range points {
		x := rnd.Float64()
		y := rnd.Float64()
		points[i] = Point{x, y}
	}
	return points
}

func normal(n int, rnd *rand.Rand) []Point {
	points := make([]Point, n)
	for i := range points {
		x := rnd.NormFloat64()
		y := rnd.NormFloat64()
		points[i] = Point{x, y}
	}
	return points
}

func grid(n int, rnd *rand.Rand) []Point {
	side := int(math.Floor(math.Sqrt(float64(n))))
	n = side * side
	points := make([]Point, 0, n)
	for y := 0; y < side; y++ {
		for x := 0; x < side; x++ {
			p := Point{float64(x), float64(y)}
			points = append(points, p)
		}
	}
	return points
}

func circle(n int, rnd *rand.Rand) []Point {
	points := make([]Point, n)
	for i := range points {
		t := float64(i) / float64(n)
		x := math.Cos(t)
		y := math.Sin(t)
		points[i] = Point{x, y}
	}
	return points
}

func shouldFail(t *testing.T, points []Point) {
	if _, err := Triangulate(points); err == nil {
		t.Fatalf("should have failed. %v", points)
	}
}

func validate(t *testing.T, points []Point) *Triangulation {
	tri, err := Triangulate(points)
	if err != nil {
		t.Fatal(err)
	}
	err = tri.Validate()
	if err != nil {
		t.Fatal(err)
	}
	return tri
}

func TestTricky(t *testing.T) {
	var points []Point
	rnd := rand.New(rand.NewSource(99))
	for len(points) < 100000 {
		x := rnd.NormFloat64() * 0.5
		y := rnd.NormFloat64() * 0.5
		points = append(points, Point{x, y})
		nx := rnd.Intn(4)
		for i := 0; i < nx; i++ {
			x = math.Nextafter(x, x+1)
		}
		ny := rnd.Intn(4)
		for i := 0; i < ny; i++ {
			y = math.Nextafter(y, y+1)
		}
		points = append(points, Point{x, y})
	}
	validate(t, points)
}

func TestCases(t *testing.T) {
	validate(t, nil)
	validate(t, []Point{{516, 661}, {369, 793}, {426, 539}, {273, 525}, {204, 694}, {747, 750}, {454, 390}})
	validate(t, []Point{{382, 302}, {382, 328}, {382, 205}, {623, 175}, {382, 188}, {382, 284}, {623, 87}, {623, 341}, {141, 227}})
	validate(t, []Point{{4, 1}, {3.7974166882130675, 2.0837249985614585}, {3.2170267516619773, 3.0210869309396715}, {2.337215067329615, 3.685489874065187}, {1.276805078389906, 3.9872025288851036}, {0.17901102978375127, 3.885476929518457}, {-0.8079039091377689, 3.3940516818407187}, {-1.550651407188842, 2.5792964886320684}, {-1.9489192990517052, 1.5512485534497125}, {-1.9489192990517057, 0.44875144655029087}, {-1.5506514071888438, -0.5792964886320653}, {-0.8079039091377715, -1.394051681840717}, {0.17901102978374794, -1.8854769295184561}, {1.276805078389902, -1.987202528885104}, {2.337215067329611, -1.6854898740651891}, {3.217026751661974, -1.021086930939675}, {3.7974166882130653, -0.08372499856146409}})
	validate(t, []Point{{0, 0}, {0, 0}, {1, 0}, {1, 1}, {1, 0}, {1, 1}})
	validate(t, []Point{{66.103648384371410884341457858681679, 68.588612471664760050771292299032211}, {146.68071346210041383528732694685459, 121.68071346210042804614204214885831}, {128.86889656046744789819058496505022, 117.26179755904141188693756703287363}, {66.103648384371439306050888262689114, 68.588612471664774261626007501035929}, {169.55213966757199273160949815064669, 146.13377653827689073295914568006992}, {126.62939224605088384123519062995911, 181.11140466039208263282489497214556}, {74.434448280233709738240577280521393, 78.630898779520691732614068314433098}, {121.11140466039205421111546456813812, 153.37060775394911615876480937004089}, {98.888595339607888945465674623847008, 186.62939224605085541952576022595167}, {52.66066896814022157968793180771172, 63.178539267712423566081270109862089}, {85.321337936280443159375863615423441, 86.357078535424832921307825017720461}, {129.61570560806461571701220236718655, 173.90180644032261625397950410842896}, {91.52240934977427855301357340067625, 162.34633135269814374623820185661316}, {137.24095128280055178038310259580612, 112.2409512828005375695283873938024}, {93.370607753949116158764809370040894, 158.88859533960791736717510502785444}, {175, 150}, {124.14213562373090837809286313131452, 184.1421356237309794323664391413331}, {96.208227592327205002220580354332924, 94.083258291328988320856296923011541}, {98.88859533960798842144868103787303, 153.37060775394905931534594856202602}, {117.98200690442070026620058342814445, 109.53561780313727069824381032958627}, {116.19447026430383118622557958588004, 108.267043413376910621082060970366}, {54.324378061245710114235407672822475, 62.306334965997713482011022279039025}, {30.886889656046740526562643935903907, 47.726179755904141188693756703287363}, {107.09511724837395263421058189123869, 101.8094380472331295095500536262989}, {38.8922619486326652804564218968153, 52.594841299088443520304281264543533}, {146.68071346210041383528732694685459, 121.68071346210039962443261174485087}, {95.857864376269077411052421666681767, 155.8578643762690205676335608586669}, {54.324378061245703008808050071820617, 62.306334965997706376583664678037167}, {137.24095128280055178038310259580612, 112.24095128280055178038310259580612}, {161.52956552860769079416058957576752, 140.440336826753821242164121940732}, {90.384294391935398493842512834817171, 166.09819355967738374602049589157104}, {113.22072967687428501903923461213708, 93.717722494332946325812372379004955}, {77.882918707497154287011653650552034, 74.870889977331813724958919920027256}, {50, 60}, {85.321337936280457370230578817427158, 86.357078535424847132162540219724178}, {41.773779312093481053125287871807814, 55.452359511808289482814871007576585}, {89.662189030622869267972419038414955, 81.153167482998867399146547541022301}, {101.44145935374857003807846922427416, 87.435444988665906862479459960013628}, {124.14213562373096522151172393932939, 155.85786437626904898934299126267433}, {172.41645518465438158273173030465841, 148.16651658265794822000316344201565}, {63.547558624186912368259072536602616, 70.9047190236165221222108812071383}, {150.64267587256094316217058803886175, 132.71415707084969426432508043944836}, {109.99999999999992894572642398998141, 190}, {128.47759065022572144698642659932375, 177.65366864730182783205236773937941}, {90, 169.99999999999994315658113919198513}, {128.47759065022574986869585700333118, 162.34633135269820058965706266462803}, {156.12047564140027589019155129790306, 131.12047564140027589019155129790306}, {90.384294391935384282987797632813454, 173.90180644032250256714178249239922}, {95.857864376268992145924130454659462, 184.1421356237308941672381479293108}, {77.882918707497140076156938448548317, 74.870889977331799514104204718023539}, {139.75578621651419553018058650195599, 124.98797731494555307563132373616099}, {130, 170}, {102.34633135269812953538348665460944, 188.47759065022569302527699619531631}, {41.773779312093481053125287871807814, 55.452359511808282377387513406574726}, {91.522409349774235920449427794665098, 177.65366864730171414521464612334967}, {27.784523897265298586489734589122236, 45.189682598176865724326489726081491}, {126.62939224605091226294462103396654, 158.88859533960797421059396583586931}, {106.09819355967735532431106548756361, 189.61570560806458729530277196317911}, {52.660668968140200263405859004706144, 63.178539267712395144371839705854654}, {74.434448280233681316531146876513958, 78.63089877952067752175935311242938}, {106.09819355967746901114878710359335, 150.38429439193538428298779763281345}, {117.65366864730172835606936132535338, 188.47759065022574986869585700333118}, {125, 100}, {38.892261948632565804473415482789278, 52.594841299088379571458062855526805}, {52.660668968140228685115289408713579, 63.17853926771241646065391250886023}, {129.61570560806461571701220236718655, 166.09819355967744058943935669958591}, {20, 40}, {117.65366864730181362119765253737569, 151.52240934977427855301357340067625}, {161.52956552860766237245115917176008, 140.440336826753821242164121940732}, {63.547558624186969211677933344617486, 70.904719023616564754775026813149452}, {127.80118910350067551462416304275393, 102.80118910350067551462416304275393}, {89.66218903062284084626298863440752, 81.153167482998853188291832339018583}, {102.34633135269824322222120827063918, 151.52240934977425013130414299666882}, {93.370607753949059315345948562026024, 181.11140466039196894598717335611582}, {113.90180644032250256714178249239922, 189.61570560806461571701220236718655}, {121.11140466039199736769660376012325, 186.62939224605094068465405143797398}, {113.90180644032258783227007370442152, 150.38429439193538428298779763281345}, {110.00000000000002842170943040400743, 150}, {165.56023782070013794509577564895153, 140.56023782070013794509577564895153}})
	validate(t, []Point{{63.547558624186912368259072536602616, 70.9047190236165221222108812071383}, {63.547558624186969211677933344617486, 70.904719023616564754775026813149452}, {66.103648384371410884341457858681679, 68.588612471664760050771292299032211}, {77.882918707497154287011653650552034, 74.870889977331813724958919920027256}, {128.47759065022572144698642659932375, 177.65366864730182783205236773937941}})
	validate(t, []Point{{54.324378061245710114235407672822475, 62.306334965997713482011022279039025}, {63.547558624186912368259072536602616, 70.9047190236165221222108812071383}, {63.547558624186969211677933344617486, 70.904719023616564754775026813149452}, {90.384294391935398493842512834817171, 166.09819355967738374602049589157104}, {90, 169.99999999999994315658113919198513}})
	shouldFail(t, []Point{{0, 0}})
	shouldFail(t, []Point{{0, 0}, {0, 0}, {0, 0}})
	shouldFail(t, []Point{{0, 0}, {1, 0}})
	shouldFail(t, []Point{{0, 0}, {1, 0}, {2, 0}, {3, 0}})
	shouldFail(t, []Point{{0, 0}, {1, 0}, {2, 0}, {3, 1}, {4, math.Inf(1)}})
	shouldFail(t, []Point{{0, 0}, {1, 0}, {2, 0}, {3, 1}, {math.Inf(1), 4}})
	shouldFail(t, []Point{{0, 0}, {1, 0}, {2, 0}, {3, 1}, {4, math.Inf(-1)}})
	shouldFail(t, []Point{{0, 0}, {1, 0}, {2, 0}, {3, 1}, {math.Inf(-1), 4}})
	shouldFail(t, []Point{{0, 0}, {1, 0}, {2, 0}, {3, 1}, {math.Inf(-1), math.NaN()}})
}

func TestUniform(t *testing.T) {
	rnd := rand.New(rand.NewSource(99))
	points := uniform(100000, rnd)
	validate(t, points)
}

func TestNormal(t *testing.T) {
	rnd := rand.New(rand.NewSource(99))
	points := normal(100000, rnd)
	validate(t, points)
}

func TestGrid(t *testing.T) {
	rnd := rand.New(rand.NewSource(99))
	points := grid(100000, rnd)
	tri := validate(t, points)

	// additional grid testing
	ts := tri.Triangles
	for i := 0; i < len(ts); i += 3 {
		p0 := points[ts[i+0]]
		p1 := points[ts[i+1]]
		p2 := points[ts[i+2]]
		a := area(p0, p1, p2)
		if a != 1 { // parallelogram area
			t.Fatal("invalid grid triangle area")
		}
	}
}

func TestCircle(t *testing.T) {
	rnd := rand.New(rand.NewSource(99))
	points := circle(10000, rnd)
	validate(t, points)
}

func BenchmarkUniform(b *testing.B) {
	rnd := rand.New(rand.NewSource(99))
	points := uniform(b.N, rnd)
	Triangulate(points)
}

func BenchmarkNormal(b *testing.B) {
	rnd := rand.New(rand.NewSource(99))
	points := normal(b.N, rnd)
	Triangulate(points)
}

func BenchmarkGrid(b *testing.B) {
	rnd := rand.New(rand.NewSource(99))
	points := grid(b.N, rnd)
	Triangulate(points)
}
