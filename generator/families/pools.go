package families

var cities = []string{
	"Lisbon", "Tallinn", "Porto", "Ljubljana", "Valencia", "Krakow", "Bergen",
	"Ghent", "Bologna", "Tbilisi", "Seville", "Gdansk", "Leipzig", "Utrecht",
	"Aarhus", "Brno", "Cluj", "Graz", "Malmo", "Riga", "Turku", "Vilnius",
	"Zadar", "Lyon", "Nantes", "Bilbao", "Salzburg", "Plovdiv", "Tartu",
	"Coimbra", "Trieste", "Kaunas", "Pecs", "Bratislava", "Uppsala", "Lille",
	"Groningen", "Innsbruck", "Varna", "Debrecen", "Oulu", "Antwerp", "Bruges",
	"Rotterdam", "Eindhoven", "Maastricht", "Hamburg", "Dresden", "Freiburg",
	"Heidelberg", "Nuremberg", "Bremen", "Hanover", "Zurich", "Basel", "Geneva",
	"Lausanne", "Lucerne", "Vienna", "Linz", "Prague", "Ostrava", "Olomouc",
	"Wroclaw", "Poznan", "Lublin", "Szeged", "Budapest", "Timisoara", "Iasi",
	"Sibiu", "Burgas", "Skopje", "Ohrid", "Belgrade", "Novi Sad", "Zagreb",
	"Rijeka", "Sarajevo", "Mostar", "Podgorica", "Tirana", "Thessaloniki",
	"Patras", "Heraklion", "Naples", "Palermo", "Turin", "Verona", "Genoa",
	"Parma", "Padua", "Marseille", "Bordeaux", "Toulouse", "Strasbourg",
	"Montpellier", "Rennes", "Grenoble", "Granada", "Malaga", "Zaragoza",
	"Cordoba", "Santander", "Oviedo", "Braga", "Faro", "Aveiro", "Galway",
	"Limerick", "Belfast", "Glasgow", "Edinburgh", "Aberdeen", "Cardiff",
	"Bristol", "Leeds", "Oslo", "Trondheim", "Stavanger", "Gothenburg", "Lund",
	"Odense", "Aalborg", "Tampere", "Espoo", "Helsinki", "Reykjavik", "Akureyri",
	"Osaka", "Kyoto", "Sapporo", "Fukuoka", "Nagoya", "Busan", "Daegu", "Incheon",
	"Taipei", "Kaohsiung", "Hanoi", "Da Nang", "Bangkok", "Chiang Mai", "Penang",
	"Ipoh", "Cebu", "Davao", "Jakarta", "Bandung", "Surabaya", "Yogyakarta",
	"Kathmandu", "Pokhara", "Colombo", "Kandy", "Mumbai", "Pune", "Jaipur",
	"Udaipur", "Kochi", "Mysore", "Chennai", "Hyderabad", "Lahore", "Karachi",
	"Muscat", "Doha", "Amman", "Beirut", "Tunis", "Marrakesh", "Fez", "Casablanca",
	"Dakar", "Accra", "Kumasi", "Lagos", "Nairobi", "Mombasa", "Arusha", "Kigali",
	"Kampala", "Lusaka", "Harare", "Windhoek", "Durban", "Cape Town", "Montreal",
	"Quebec City", "Halifax", "Calgary", "Winnipeg", "Vancouver", "Seattle",
	"Portland", "Denver", "Memphis", "Boston", "Chicago", "Detroit", "Phoenix",
	"Tucson", "Oaxaca", "Puebla", "Merida", "Guadalajara", "Havana", "Bogota",
	"Medellin", "Cali", "Quito", "Cuenca", "Lima", "Cusco", "Arequipa", "La Paz",
	"Santiago", "Valparaiso", "Mendoza", "Rosario", "Montevideo", "Asuncion",
	"Recife", "Curitiba", "Manaus", "Auckland", "Wellington", "Hobart", "Perth",
	"Adelaide", "Brisbane", "Darwin",
}

// Pool sizes bound MaxCases: Builder.Pick never reuses a value within a
// pack, and one round of the fifteen default families draws at most 22
// cities, 9 projects, 5 milestones, 3 statuses, 2 firms, 2 domains and one
// nickname. MaxCases (80) is six rounds; TestPoolsFillMaxCases holds that.
// Values must also stay apart: no value may contain another value from any
// pool, so one case's gold can never ground or decoy another's.

var firstNames = []string{
	"Ines", "Tobias", "Marisol", "Henrik", "Aiyana", "Dmitri", "Noor", "Callum",
	"Yuki", "Bastian", "Leilani", "Oskar", "Priya", "Teodor", "Ama", "Rasmus",
	"Zofia", "Emeka", "Liesel", "Joaquin", "Saoirse", "Kenji", "Farida", "Mateo",
	"Ingrid", "Tariq", "Ottilie", "Ravi", "Signe", "Dario", "Hana", "Lorcan",
	"Mirela", "Anselm", "Keziah", "Bram", "Yasmin", "Florin", "Thandiwe", "Ilan",
	"Wren", "Casimir", "Odalys", "Pavel", "Nkiru", "Soren", "Ximena", "Halvard",
	"Annika", "Dagny", "Kofi", "Rosalind", "Adaeze", "Bjorn", "Celestine",
	"Darius", "Elodie", "Fintan", "Gisela", "Hamish", "Ifeoma", "Jarrah",
	"Katarina", "Leopold", "Magnus", "Nadia", "Orla", "Philippa", "Quentin",
	"Rhiannon", "Stellan", "Tamsin", "Ulrich", "Valentina", "Wilhelmina",
	"Xavier", "Yara", "Zoltan", "Anouk", "Benedikt", "Cosima", "Desmond",
	"Esperanza", "Fabian", "Greta", "Horatio", "Imogen", "Jovan", "Klaus",
	"Ludmila", "Mihail", "Nikolai", "Ophelia", "Piotr", "Radu", "Svetlana",
	"Thaddeus", "Vesna", "Wolfgang", "Yevgenia", "Zuzana", "Agnieszka", "Bartosz",
	"Cormac", "Dorotea", "Eamon", "Fenna", "Gunnar", "Hedda", "Ignatius",
	"Jolanda", "Kasimir", "Lavinia", "Marek", "Niamh", "Olafur", "Petra",
	"Ragnhild", "Sebastiano", "Tomasz", "Ursula", "Viktor", "Wiebke", "Yannick",
	"Zelda", "Aurelio", "Beatrix", "Cyprian", "Delphine", "Emrys", "Frieda",
}

var events = []string{
	"team offsite", "book club retreat", "alumni reunion", "product summit",
	"family reunion", "choir tour", "wedding rehearsal", "charity gala",
	"hackathon", "quiz night", "garden party", "board game night",
	"cooking class", "film screening", "photo walk", "baby shower",
	"retirement party", "housewarming", "chess tournament", "pottery workshop",
	"wine tasting", "farewell dinner", "science fair", "yoga retreat",
	"trivia final", "sailing regatta",
}

var filler = [][2]string{
	{"Remind me to water the plants on Thursday.", "I'll remind you Thursday."},
	{"Watched a documentary about octopuses last night.", "They're fascinating."},
	{"Thinking about learning to bake sourdough.", "Start with a simple starter."},
	{"It's been unusually windy this week.", "Windy weeks can be draining."},
	{"Need to renew my passport sometime this year.", "Worth starting early."},
	{"My bike chain keeps slipping on hills.", "It might need a new chain or a tune-up."},
	{"Finally finished the jigsaw puzzle, all thousand pieces.", "That's a satisfying feeling."},
	{"The neighbours are redoing their kitchen, so it's loud all day.", "Noise-cancelling headphones might help."},
	{"I keep forgetting where I put my reading glasses.", "A fixed spot by the door can help."},
	{"Started journaling before bed.", "That's a nice way to wind down."},
	{"My laptop fan has been really loud lately.", "Cleaning the vents often fixes that."},
	{"Got a parking ticket outside the library, annoying.", "Sorry, that's frustrating."},
	{"Thinking of repainting the hallway a lighter colour.", "Lighter colours make hallways feel bigger."},
	{"The new season of that detective show is out.", "Enjoy the binge."},
	{"I signed up for a beginner pottery class online.", "Sounds like a relaxing hobby."},
	{"My phone screen cracked again.", "A sturdier case might be worth it."},
	{"We're thinking about getting a rescue dog.", "Rescue dogs can be wonderful companions."},
	{"Spent the whole weekend sorting old photos.", "That can bring back a lot of memories."},
	{"The printer at work jammed three times today.", "Printers have a talent for bad timing."},
	{"Trying to drink more water during the day.", "A marked bottle can help you keep track."},
	{"I should really back up my files this week.", "Good idea, before something goes wrong."},
	{"Had a strange dream about missing a train.", "Those dreams are surprisingly common."},
	{"The streetlight outside my window keeps flickering.", "Worth reporting it to the council."},
	{"I'm trying to read one book a month this year.", "That's a very doable goal."},
	{"Our team switched to a new chat app at work.", "Switching tools always takes some getting used to."},
	{"My houseplant finally put out a new leaf.", "Nice, it must like its spot."},
	{"Forgot my umbrella and got soaked on the way home.", "Hope you dried off quickly."},
	{"Looking for a good podcast for my commute.", "Do you prefer stories or interviews?"},
	{"The gym was packed this morning.", "Early mornings get busy after the holidays."},
	{"I've been meaning to call my grandmother.", "She'd probably love to hear from you."},
}

// surnames pair with firstNames in Builder.Person; all are plain
// capitalised words, the shape the agents' full-name pattern reads.
var surnames = []string{
	"Varga", "Moreau", "Haas", "Lindqvist", "Okafor", "Brandt", "Castillo",
	"Novak", "Petrov", "Sato", "Kowalski", "Fischer", "Romero", "Duarte",
	"Ferreira", "Nakamura", "Larsen", "Holm", "Dubois", "Weber", "Ibarra",
	"Quinlan", "Marsh", "Whitfield", "Achebe", "Mendes", "Rossi", "Bianchi",
	"Keller", "Vogel", "Sandoval", "Pereira", "Coelho", "Takahashi", "Yilmaz",
	"Demir", "Kaya", "Horvath", "Szabo", "Toth", "Nagy", "Dvorak", "Svoboda",
	"Jensen", "Nilsen", "Berg", "Strand", "Aalto", "Virtanen", "Kask", "Ozols",
	"Kalnins", "Popescu", "Ionescu", "Dimitrov", "Georgiev", "Papadakis", "Costa",
	"Russo", "Gallo", "Conti", "Ortega", "Navarro", "Molina", "Delgado",
	"Herrera", "Medina", "Vargas", "Rios", "Fuentes", "Acosta", "Pacheco",
	"Salazar", "Bauer", "Hartmann", "Kruger", "Lehmann", "Schulz", "Wagner",
	"Zimmer", "Abara", "Mensah", "Owusu", "Boateng", "Adeyemi", "Nwosu", "Kamau",
	"Otieno", "Mwangi", "Banda", "Phiri", "Tembo", "Kim", "Park", "Choi", "Lim",
	"Tan", "Wong", "Chua", "Ong", "Reyes", "Santos", "Cruz", "Bautista",
}

// nicknames are what colleagues call someone; NicknameEmail changes an
// address under the nickname alone.
var nicknames = []string{
	"Kit", "Doc", "Ace", "Bee", "Sunny", "Moose", "Pip", "Biscuit", "Sparky",
	"Duke", "Scout", "Rocket", "Bear", "Taz", "Jinx", "Ziggy", "Boots", "Pepper",
	"Maverick", "Tink", "Dash", "Fizz", "Sprout", "Chip", "Bo", "Rusty", "Mojo",
	"Noodle", "Peaches", "Gizmo", "Bunny", "Coco", "Dottie", "Flick", "Ginger",
	"Honey", "Juno", "Kiki", "Lulu", "Mittens", "Nugget", "Pebbles", "Rascal", "Skip",
	"Tiger", "Waffles", "Zippy", "Button", "Cricket", "Dimples", "Freckles",
	"Bramble", "Clover", "Domino", "Fudge", "Goose", "Hobbit", "Jellybean", "Kipper",
	"Marmalade", "Nibbles", "Pickle", "Rooster", "Sherlock", "Toffee", "Wiggles", "Zest",
	"Acorn", "Beetle", "Cobbler", "Dumpling", "Fable", "Gadget", "Hopper", "Jester",
	"Kettle", "Lemon", "Mango", "Nutmeg", "Otter", "Puddle", "Quokka", "Ripple",
	"Sardine", "Tadpole", "Whisker", "Yeti",
}

// shopNames are coined first words of shop names (alias-list). Coined, so a
// shop never matches a real word in another pool or the filler.
var shopNames = []string{
	"Crumbwell", "Tallowmere", "Fernhollow", "Brackenly", "Quillmoor", "Marrowby", "Pellingham",
	"Rookwell", "Tansybrook", "Wexcombe", "Hollinshed", "Larchwick", "Mossington", "Penderleigh",
	"Saffronby", "Thornquist", "Umberly", "Varnley", "Wickerham", "Yarrowdene", "Bellcote",
	"Copperfold", "Dunmarsh", "Eldermoor", "Fallowgate", "Gorsewick", "Heathcombe", "Ivyholt",
	"Juniperly", "Kestrelby", "Lindenwold", "Meadowick", "Nettlebury", "Oakhaven", "Primshaw",
	"Quarrowby", "Ravelston", "Sorrelby", "Teasleworth", "Underwick", "Vetchley", "Willowmere",
	"Brambleton", "Cloverdine", "Dappleford", "Ambleworth", "Briarholt", "Candlemoor",
	"Dewbury", "Elmscote", "Foxbury", "Glenmarrow", "Hartlecombe", "Inglewick", "Jasperly",
	"Kingsholt", "Lowenbury", "Mistlecote", "Northwick", "Oxlade", "Pennyrose", "Quenby",
	"Redmarsh", "Stellington", "Thistlecote", "Upperleigh", "Vernholt", "Wrenfield", "Yewbury",
	"Zellby",
}

// projects are coined names, so no project is a real word another pool or
// the filler could contain.
var projects = []string{
	"Aldoria", "Brevant", "Calimor", "Dovrane", "Esterly", "Falkira", "Gorvane",
	"Halcyra", "Istrel", "Kelvane", "Lunaris", "Morvett", "Nexara", "Ostrin",
	"Pellara", "Quorra", "Ravelin", "Solvane", "Ulvane", "Varnis", "Wexley",
	"Zelvane", "Ambrel", "Bastrel", "Corvin", "Delvar", "Emberly", "Fenwick",
	"Galdor", "Hesper", "Ivoran", "Kalvex", "Lorvane", "Mirelle", "Norrin",
	"Orvell", "Pravin", "Quentro", "Rovane", "Selmar", "Ulmora", "Vespera",
	"Wrenfold", "Xylar", "Yoravel", "Zanthe", "Arcadel", "Brimley", "Caskel",
	"Dunmore", "Elstow", "Fairlow", "Glenvar", "Holloway", "Jesmond", "Kirrin",
	"Lowmere", "Marrick", "Nettley", "Oakhurst", "Pemberly", "Rookwood",
	"Thornbury", "Upworth", "Valemont", "Westbrook", "Yewdale", "Zennor",
	"Ashgrove", "Birchmont", "Coldwell", "Everdene", "Foxmere", "Greywell",
	"Hartcombe", "Kingsmere", "Larchfield", "Millbrook", "Northwold", "Orchardly",
	"Pinecrest", "Redfern", "Silverlode", "Tidewater", "Underhill", "Wexmoor",
}

// domains are coined work-mail domains.
var domains = []string{
	"lumenfield.co", "corbelhq.com", "fernway.io", "quayside.net", "ridgeview.co",
	"pinemark.io", "halden.co", "brightwater.org", "ostara.io", "velvetine.co",
	"marlowe-labs.com", "tidecraft.io", "ashbourne.co", "cobaltpeak.com",
	"greystone.io", "lanternhq.co", "meridianworks.net", "oakenshaw.co",
	"pebblebay.io", "quillstone.com", "rivermint.co", "saltaire.io",
	"thistledown.org", "umberline.co", "vesperhq.com", "willowmere.io",
	"yellowpine.co", "zephyrlane.com", "amberfield.io", "birchline.co",
}

// firms hold contracts in OwnerFirm.
var firms = []string{
	"Brightline Catering", "Copperleaf Hospitality", "Saltmarsh Kitchens",
	"Larder Lane Foods", "Greenhouse Table", "Silver Spoon Events",
	"Harvest Moon Catering", "Bluebell Banquets", "Fig Tree Kitchens",
	"Olive Grove Catering", "Ember Hall Events", "Linden Street Catering",
	"Marigold Feasts", "Cedar Plate Company", "Sunday Supper Club",
	"Wildflower Catering", "Quartz Kitchen Collective", "Northstar Banqueting",
	"Tidepool Catering", "Hearthside Events", "Crimson Ladle",
	"Pepperpot Catering", "Willow Table Events", "Ironbark Kitchens",
	"Goldleaf Dining", "Riverside Spread", "Cloverfield Catering",
	"Amberline Events", "Stoneoven Catering", "Blue Door Kitchens",
	"Palisade Hospitality", "Meadowlark Catering", "Foxglove Feasts",
	"Lantern House Catering", "Seaglass Events", "Orchard Row Kitchens",
	"Peppercorn Hall", "Firefly Catering", "Driftwood Dining",
}

// statuses are multi-word states with no negation, past marker or "and",
// so every one grades as one phrase.
var statuses = []string{
	"awaiting final sign-off", "cleared for printing",
	"paused pending budget review", "with the legal team",
	"waiting on vendor quotes", "scheduled for spring", "approved in principle",
	"in final testing", "out for review", "ready to ship",
	"blocked on procurement", "under budget review", "parked until autumn",
	"live in pilot", "queued for the board", "with the design team",
	"awaiting customer feedback", "signed off by finance", "in contract talks",
	"on the launch calendar", "pending safety check", "waiting on permits",
	"drafted for comment", "in user testing", "booked for delivery",
	"under vendor evaluation", "ready for handover", "held for legal review",
	"in the approval queue", "waiting on translation", "set for soft launch",
	"slotted for next quarter", "in final edits", "with procurement",
	"awaiting site survey", "cleared by compliance", "in price negotiation",
	"ready for sign-off", "scheduled for rollout", "back in drafting",
}

// milestones are the dated items of one schedule in EventOrder.
var milestones = []string{
	"dress fitting", "venue walkthrough", "cake tasting", "rehearsal dinner",
	"engagement shoot", "florist meeting", "menu trial", "guest list freeze",
	"seating plan review", "band audition", "invitation printing", "hair trial",
	"suit alterations", "license appointment", "ring pickup", "welcome brunch",
	"playlist review", "lighting check", "transport booking",
	"speech run-through", "budget check-in", "site visit", "kickoff call",
	"design review", "vendor demo", "contract signing", "pilot launch",
	"training session", "board briefing", "press preview", "user interviews",
	"beta freeze", "load test", "security audit", "price review", "stock count",
	"staff briefing", "photo selection", "catalog proof", "showroom visit",
	"furniture delivery", "paint sampling", "plumber visit", "electrician visit",
	"tile selection", "kitchen fitting", "carpet install", "final inspection",
	"key handover", "moving day", "garage clearout", "school tour",
	"parent evening", "sports day", "piano recital", "uniform fitting",
	"dentist checkup", "eye test", "passport appointment",
}

// threads are pieces of work a person holds or a status describes.
var threads = []string{
	"vendor audit", "office move", "brand refresh", "grant application",
	"hiring drive", "pricing update", "website rebuild", "warehouse lease",
	"insurance renewal", "client onboarding", "data migration",
	"partner contract", "safety review", "fleet upgrade", "supplier switch",
	"annual report", "kitchen remodel", "garden redesign", "loan refinance",
	"visa paperwork",
}

// rosters are the lists MembershipFold changes.
var rosters = []string{
	"carpool rota", "book club roster", "volunteer list", "party guest roster",
	"relay team", "study group", "choir roster", "snack rota", "cleaning rota",
	"board game group",
}

// schedules hold EventOrder's milestones.
var schedules = []string{"wedding plans", "renovation plan", "launch plan", "move plan", "spring project plan", "garden project"}
