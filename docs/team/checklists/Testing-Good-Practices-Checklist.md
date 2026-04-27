# Good Practices Checklist - Testing Edition

Αυτή η σύνοψη αποτελεί το απόλυτο **Checklist Καλών Πρακτικών** για τη συγγραφή, διαχείριση και βελτιστοποίηση των Tests στην Go, βασισμένη στην επίσημη τεκμηρίωση του πακέτου `testing`.

Χρησιμοποίησε αυτή τη λίστα για να βεβαιωθείς ότι το testing suite του project σου (είτε είναι CLI, είτε Web, είτε βιβλιοθήκη) είναι ιδιωματικό, αποδοτικό και πλήρως εκμεταλλεύσιμο.

---

## 1. Δομή & Οργάνωση (Structure & Organization)
- [ ] **Ονομασία & Τοποθεσία:** Ονομάζω τα tests ως `func TestXxx(t *testing.T)` (όπου το `Xxx` **δεν** ξεκινά με πεζό γράμμα) και τα αποθηκεύω αυστηρά σε αρχεία `*_test.go` ώστε να μην περιλαμβάνονται στα κανονικά builds.
- [ ] **Black-box vs White-box:** Διαχωρίζω σωστά τα tests: Σε white-box tests (ίδιο package) καλούνται και τα unexported σύμβολα, ενώ σε black-box tests (`package name_test`) γίνεται κανονικό import και ελέγχεται αποκλειστικά το public API.
- [ ] **Σήμανση Αποτυχιών:** Χρησιμοποιώ τα `t.Error / t.Errorf` (ή `t.Fail()`) για να σηματοδοτήσω ένα σφάλμα επιτρέποντας στο test να συνεχίσει, και `t.Fatal / t.Fatalf` μόνο όταν η συνέχεια δεν έχει νόημα.

## 2. Subtests & Ιεραρχία (Subtests & Hierarchy)
- [ ] **t.Run & b.Run:** Χρησιμοποιώ το `t.Run(name, func)` για ιεραρχικά tests (π.χ. Table-Driven) και το `b.Run()` για sub-benchmarks, αντί να δημιουργώ δυναμικά ονόματα. Έτσι εκμεταλλεύομαι ιεραρχικά τα flags επιλογής (π.χ. `-run Foo/A=1`, `-bench` ή `-fuzz`).
- [ ] **Παράλληλη Εκτέλεση:** Καθιστώ τα subtests παράλληλα καλώντας `t.Parallel()` μέσα στο `t.Run()`. Εκμεταλλεύομαι το enclosing `t.Run` για ελεγχόμενο parallelism και ομαλό cleanup μετά την ολοκλήρωσή τους.

## 3. Benchmarks (Μετρήσεις Απόδοσης)
- [ ] **Ονομασία & Βρόχος:** Ονομάζω τα benchmarks ως `func BenchmarkXxx(b *testing.B)` και χρησιμοποιώ το σύγχρονο `for b.Loop() { ... }` (ή το `for i := 0; i < b.N; i++` για παλαιότερες εκδόσεις).
- [ ] **Ακριβό Setup:** Αν το benchmark απαιτεί "βαρύ" setup (π.χ. δημιουργία δεδομένων), το κάνω εκτός του βρόχου, ή καλώ ρητά την `b.ResetTimer()` ακριβώς πριν ξεκινήσει το loop για να μην αλλοιωθεί η μέτρηση του χρόνου.
- [ ] **Μετρήσεις Μνήμης & Μετρικές:** Χρησιμοποιώ `b.ReportAllocs()` (ή τη συνάρτηση `testing.AllocsPerRun`) για καταγραφή δεσμεύσεων μνήμης, και την `b.ReportMetric()` για custom μετρικές (π.χ. items/op).
- [ ] **Παράλληλα Benchmarks:** Για μέτρηση απόδοσης σε συνθήκες ταυτόχρονης εκτέλεσης (concurrency), χρησιμοποιώ την `b.RunParallel()` και προσαρμόζω αν χρειάζεται με την `b.SetParallelism()`.
- [ ] **Εργαλεία Σύγκρισης:** Χρησιμοποιώ εξωτερικά εργαλεία όπως το `benchstat` (αντί για απλή οπτική σύγκριση) για την παραγωγή στατιστικά έγκυρων συμπερασμάτων μεταξύ διαφορετικών εκτελέσεων κώδικα.

## 4. Fuzz Testing (Έλεγχος με Τυχαία Δεδομένα)
- [ ] **Δομή & Seed Corpus:** Δημιουργώ Fuzz tests με τη μορφή `func FuzzXxx(f *testing.F)` και παρέχω αρχικό seed corpus μέσω της `f.Add(...)` ή αποθηκεύοντας αρχεία στον φάκελο `testdata/fuzz/<FuzzName>/`.
- [ ] **Target Function:** Ορίζω τον κώδικα ελέγχου μέσα στο `f.Fuzz(func(t *testing.T, inputs...) { ... })`, φροντίζοντας οι τύποι εισόδου να ταιριάζουν με αυτούς του seed corpus.
- [ ] **Ασφαλής Απόρριψη (Skip):** Χρησιμοποιώ την `t.Skip()` μέσα στο Fuzz target όταν το τυχαία παραγόμενο input δεν είναι έγκυρο (invalid) αλλά δεν αποτελεί failure. Σηματοδοτώ bugs με `t.Fail`, `t.Error`, `t.Fatal` ή panics.
- [ ] **Εκτέλεση & Regression:** Τρέχω το fuzzing με `go test -fuzz=FuzzName`. Αν το fuzzing βρει bug, κάνω commit το παραγόμενο failing input (από τον `testdata/`) ώστε να λειτουργεί πλέον ως μόνιμο regression test.

## 5. Helpers, Setup & Teardown
- [ ] **Χρήση t.Helper():** Καλώ ρητά την `t.Helper()` (ή `b.Helper()`) στην αρχή των δικών μου βοηθητικών συναρτήσεων, ώστε αν υπάρξει σφάλμα, το log να δείξει τη γραμμή του αρχικού test και όχι τη γραμμή του helper.
- [ ] **Ασφαλές Cleanup:** Χρησιμοποιώ την `t.Cleanup()` (ή `b.Cleanup()`) αντί για απλή `defer` όταν θέλω να εγγυηθώ ένα ντετερμινιστικό teardown πόρων στο τέλος της εκτέλεσης του test ή subtest.
- [ ] **Global Setup (TestMain):** Αν απαιτείται global setup/teardown (για ολόκληρο το package), χρησιμοποιώ την `func TestMain(m *testing.M)`. Από την Go 1.15+, το wrapper διαχειρίζεται αυτόματα το exit (παλαιότερα απαιτούσε ρητά `os.Exit(m.Run())`). Αν περνάω custom flags, καλώ το `flag.Parse()` χειροκίνητα.
- [ ] **Προσωρινά Αρχεία (TempDir):** Για tests που διαβάζουν/γράφουν αρχεία, χρησιμοποιώ την `t.TempDir()` (ή `b.TempDir()` / `ArtifactDir()`), η οποία δημιουργεί προσωρινούς φακέλους που διαγράφονται αυτόματα με ασφάλεια στο τέλος του test.

## 6. Executable Documentation (Examples)
- [ ] **Γραφή & Ονοματολογία:** Γράφω παραδείγματα (examples) χρησιμοποιώντας τις αυστηρές συμβάσεις ονοματολογίας: `Example()` (για το package), `ExampleF()` (για function), `ExampleT()` (για type) και `ExampleT_M()` (για method). Για πολλαπλά examples στο ίδιο στοιχείο προσθέτω πεζό suffix (π.χ. `ExampleF_suffix`).
- [ ] **Αναμενόμενη Έξοδος:** Έχω προσθέσει το ειδικό σχόλιο `// Output:` (ή `// Unordered output:` αν η σειρά εκτύπωσης δεν είναι εγγυημένη) στο τέλος του example, ώστε το `go test` να το εκτελεί και να επαληθεύει την ορθότητά του, καθιστώντας το εκτελέσιμο documentation.

## 7. Flags & Εκτέλεση (Context)
- [ ] **Short Mode:** Χρησιμοποιώ τον έλεγχο `if testing.Short() { t.Skip("...") }` για να μπορώ να παραλείπω πολύ βαριά ή χρονοβόρα tests όταν ο χρήστης τρέχει τα tests με το flag `-short`.
- [ ] **Verbose Logging:** Χρησιμοποιώ το `testing.Verbose()` για να ελέγξω αν πρέπει να εκτυπώσω περισσότερες, αναλυτικές πληροφορίες (debugging logs) όταν τα tests τρέχουν με το flag `-v`.
- [ ] **Testing Context:** Χρησιμοποιώ την `testing.Testing()` αν θέλω ο κώδικάς μου να γνωρίζει (σε runtime) αν αυτή τη στιγμή εκτελείται μέσα σε `go test` περιβάλλον ή όχι.