# Git Workflow & Version Control

Η σωστή χρήση του Git είναι εξίσου σημαντική με τη συγγραφή καθαρού κώδικα. Ένα καθαρό ιστορικό (commit history) διευκολύνει την ομαδική συνεργασία, το debugging και λειτουργεί ως ζωντανή τεκμηρίωση (documentation) της εξέλιξης του project.

Εδώ συγκεντρώνονται οι βέλτιστες πρακτικές της βιομηχανίας για το Git Workflow.

## 1. Στρατηγικές Διακλάδωσης (Branching Strategies)

Η επιλογή της σωστής στρατηγικής εξαρτάται από το πώς γίνεται το release του λογισμικού. Μην δουλεύεις ποτέ απευθείας στο κύριο branch (`main` ή `master`).

### Α. GitHub Flow (Trunk-based Development)
Ιδανικό για CI/CD, web εφαρμογές, APIs και ευέλικτες ομάδες.
- Το **`main`** branch είναι *πάντα* λειτουργικό και έτοιμο για παραγωγή.
- Δημιουργείς βραχύβια branches (π.χ. `feat/user-auth`, `fix/login-bug`) απευθείας από το `main`.
  ```bash
  # Ενημέρωσε το τοπικό σου main και δημιούργησε ένα νέο branch
  git checkout main
  git pull origin main
  git checkout -b feat/add-login-page
  ```
- Ανοίγεις Pull Request, περνάει Code Review, και γίνεται merge πίσω στο `main`.

### Β. Git Flow (Το Παραδοσιακό Μοντέλο)
Ιδανικό για λογισμικό με αυστηρές εκδόσεις (π.χ. Desktop Apps, Games).
- **`main`**: Περιέχει μόνο κώδικα που έχει γίνει επίσημο release.
- **`develop`**: Το "ενεργό" branch (integration) όπου συγχωνεύονται τα νέα χαρακτηριστικά.
- Χρησιμοποιεί αυστηρά υποστηρικτικά branches: `feature/` (νέος κώδικας), `release/` (προετοιμασία έκδοσης), και `hotfix/` (επείγουσες διορθώσεις απευθείας στο `main`).
  ```bash
  # Δημιουργία ενός feature branch από το develop
  git checkout develop
  git pull origin develop
  git checkout -b feature/new-inventory-system
  ```

## 2. Ατομικά Commits (Atomic Commits)

Είναι πολύ δελεαστικό να τρέξεις `git add .` και να τα βάλεις όλα σε ένα τεράστιο commit, αλλά αυτό είναι κακή πρακτική. Ένα commit πρέπει να αφορά **μία και μόνο λογική αλλαγή**.

Μην ομαδοποιείς την αλλαγή στο README, μαζί με τη διόρθωση ενός bug, μαζί με την προσθήκη ενός νέου feature. Σπάσε τα σε 3 ξεχωριστά commits.

- **Πρακτικά:** Χρησιμοποίησε πάντα `git status` και `git diff` πριν κάνεις add. Ιδανικά, μάθε το `git add -p` (patch) για να προσθέτεις συγκεκριμένα κομμάτια κώδικα (hunks) τη φορά.
  ```bash
  # Δες τις αλλαγές σου
  git status
  git diff
  
  # Πρόσθεσε διαδραστικά μόνο τις αλλαγές που αφορούν το ένα feature
  git add -p path/to/your/file.go
  ```

## 3. Συμβατικά Μηνύματα (Conventional Commits)

Ένα επαγγελματικό μήνυμα commit λέει μια ιστορία. Χρησιμοποίησε το πρότυπο **Conventional Commits**, το οποίο βάζει ένα πρόθεμα (type) πριν από το μήνυμα:

Για τον πλήρη οδηγό, τις βέλτιστες πρακτικές και τα anti-patterns, συμβουλέψου το έγγραφο **[Conventional Commits Guide](./03.%20Conventional-Commits.md)**.

## 4. Merge vs. Rebase (Διατήρηση Καθαρού Ιστορικού)

Καθώς δουλεύεις στο branch σου, το `main` branch ίσως προχωρήσει (άλλοι developers κάνουν merge δικό τους κώδικα). Πρέπει να φέρεις αυτές τις νέες αλλαγές στο branch σου.

### `git merge main`
- **Τι κάνει:** Δημιουργεί ένα νέο "Merge Commit" που ενώνει το ιστορικό του `main` με το branch σου.
- **Πλεονέκτημα:** Απόλυτα ασφαλές, διατηρεί το ακριβές ιστορικό (non-destructive).
- **Μειονέκτημα:** Δημιουργεί περιττά merge commits ("θόρυβος") που κάνουν το γράφημα (history graph) να μοιάζει με "ιστό αράχνης".
  ```bash
  # Πήγαινε στο branch σου
  git checkout feat/add-login-page
  # Φέρε τις τελευταίες αλλαγές και κάνε merge το main στο branch σου
  git fetch origin
  git merge origin/main
  ```

### `git rebase main`
- **Τι κάνει:** Αποκόπτει τα δικά σου commits, ενημερώνει το branch σου με τον πιο πρόσφατο κώδικα του `main`, και εφαρμόζει (replays) τα δικά σου commits *πάνω* από αυτό.
- **Πλεονέκτημα:** Δημιουργεί ένα απόλυτα **γραμμικό, καθαρό ιστορικό** (Linear History).
- **Μειονέκτημα:** Αλλάζει το ιστορικό (rewrites history).
  ```bash
  git checkout feat/add-login-page
  git fetch origin
  git rebase origin/main
  ```

> ⚠️ **Ο Χρυσός Κανόνας του Rebase:** Μην κάνεις ΠΟΤΕ rebase σε branches που έχουν γίνει push σε remote server και δουλεύουν άλλοι σε αυτά! Το rebase είναι ιδανικό για να κρατάς το *τοπικό* (local) και *προσωπικό* σου branch ενημερωμένο.

## 5. Interactive Rebase (`git rebase -i`) & Squashing

Πριν ανοίξεις ένα Pull Request, ίσως το τοπικό σου branch να είναι γεμάτο από "WIP" (Work In Progress) ή "fix typo" commits. Το Interactive Rebase σου επιτρέπει να καθαρίσεις αυτό το χάος.

Εκτέλεσε την παρακάτω εντολή για να διαχειριστείς τα 3 τελευταία commits σου:
```bash
git rebase -i HEAD~3
```
Αυτό θα ανοίξει τον editor σου με τα commits (από το παλαιότερο στο νεότερο):
```text
pick 3a1b2c3 feat: add user auth logic
pick 8f9g0h1 fix: typo in auth
pick 4j5k6l7 refactor: cleanup auth errors
```
Αλλάζοντας τη λέξη `pick` σε `squash` (ή `s`), το Git θα συγχωνεύσει τα commits στο προηγούμενο, αφήνοντάς σου την ευκαιρία να γράψεις ένα καθαρό, ενιαίο μήνυμα:
```text
pick 3a1b2c3 feat: add user auth logic
squash 8f9g0h1 fix: typo in auth
squash 4j5k6l7 refactor: cleanup auth errors
```
Το αποτέλεσμα; Ένα μόνο, όμορφο και συμπαγές (Atomic) commit έτοιμο για PR!

## 6. Διόρθωση Παλαιών Commits (Amend & Reword)

Αν έκανες ένα ορθογραφικό λάθος σε ένα μήνυμα commit ή ξέχασες να προσθέσεις ένα αρχείο, δεν χρειάζεται να κάνεις νέο commit τύπου "fix typo". Μπορείς να το διορθώσεις διατηρώντας το ιστορικό σου καθαρό.

**Α. Διόρθωση του αμέσως προηγούμενου commit (`--amend`):**
Αν το λάθος έγινε στο *τελευταίο* σου commit και δεν το έχεις κάνει ακόμα push:
```bash
# Διορθώνει μόνο το μήνυμα του τελευταίου commit
git commit --amend -m "feat: add user auth logic with correct spelling"

# Αν ξέχασες να κάνεις add ένα αρχείο, πρόσθεσέ το στο ίδιο commit!
git add forgotten-file.go
git commit --amend --no-edit
```

**Β. Διόρθωση παλαιότερων commits (`rebase -i` και `reword`):**
Αν το ορθογραφικό λάθος βρίσκεται 3 commits πίσω, χρησιμοποιούμε ξανά το Interactive Rebase.
```bash
git rebase -i HEAD~3
```
Στον editor που θα ανοίξει, βρες το commit με το λάθος και άλλαξε τη λέξη `pick` σε `reword` (ή `r`):
```text
reword 8f9g0h1 fix: typoo in auth logic
pick 4j5k6l7 refactor: cleanup auth errors
```
Μόλις κλείσεις (αποθηκεύσεις) τον editor, το Git θα σταματήσει στο συγκεκριμένο commit και θα σου ανοίξει ένα νέο παράθυρο για να γράψεις το σωστό μήνυμα.

> ⚠️ **Προσοχή:** Αν έχεις ήδη κάνει push αυτά τα commits στο remote (π.χ. GitHub/Gitea), αφού τα διορθώσεις τοπικά θα αλλάξει το hash τους. Θα χρειαστεί να κάνεις force push (`git push origin <branch-name> --force`). Χρησιμοποίησέ το **μόνο** στο δικό σου branch, ποτέ στο `main`!

## 7. Ακύρωση Τοπικών Commits (Git Reset)

Αν έκανες "κατά λάθος" κάποια τοπικά commits και συνειδητοποίησες ότι το γράφημά σου (linear graph) πρόκειται να "σπάσει", μπορείς να τα ακυρώσεις **χωρίς να χάσεις τον κώδικά σου**, ώστε να τα ξανακάνεις commit σωστά.

- **`git reset --soft HEAD~1`**: Ακυρώνει το τελευταίο commit (ή `HEAD~2` για τα δύο τελευταία), αλλά αφήνει όλες τις αλλαγές των αρχείων σου στο "Staging Area" (πράσινα).
- **`git reset --mixed HEAD~1`**: (Προεπιλογή) Ακυρώνει το commit και βγάζει τα αρχεία από το staging (κόκκινα), έτοιμα για νέα επεξεργασία.

```bash
# Ακύρωση του τελευταίου commit κρατώντας τις αλλαγές έτοιμες
git reset --soft HEAD~1

# Τώρα μπορείς να κάνεις pull/rebase για να συγχρονίσεις το δέντρο σου
# και μετά να δημιουργήσεις ένα νέο, σωστό commit.
```

## 8. Ενημέρωση Τοπικού Branch (Pull Rebase)

Έχεις κάνει κάποια δικά σου commits τοπικά, αλλά εν τω μεταξύ προστέθηκε νέος κώδικας στο remote repository. Αν κάνεις ένα απλό `git pull`, το Git θα δημιουργήσει ένα περιττό "Merge Commit", καταστρέφοντας το γραμμικό ιστορικό.

Για να "αποκρύψεις" προσωρινά τα commits σου, να κατεβάσεις τις πρόσφατες αλλαγές, και στη συνέχεια να επανεμφανίσεις τα δικά σου commits *πάνω* από αυτές, χρησιμοποιείς το **Pull with Rebase**.

```bash
# Κατεβάζει τον νέο κώδικα και "ξαναπαίζει" τα δικά σου commits από πάνω του
git pull --rebase origin main

# Το ιστορικό σου είναι πλέον απόλυτα γραμμικό και έτοιμο για push!
git push origin my-branch
```

## 9. Προσωρινή Αποθήκευση Δουλειάς (Git Stash)

Συχνά δουλεύεις σε ένα feature και ξαφνικά πρέπει να αλλάξεις branch για να διορθώσεις ένα επείγον bug (hotfix), αλλά ο κώδικάς σου δεν είναι έτοιμος για commit.

- **`git stash push -m "μήνυμα"`:** Αποθηκεύει προσωρινά τις μη ολοκληρωμένες αλλαγές σου και επαναφέρει το branch στην τελευταία καθαρή του κατάσταση (HEAD).
- **`git stash list`:** Εμφανίζει όλες τις αποθηκευμένες "στοίβες" (stashes) σου.
- **`git stash pop`:** Εφαρμόζει ξανά την τελευταία αποθηκευμένη δουλειά στο branch που βρίσκεσαι τώρα και τη διαγράφει από τη λίστα του stash.
- **`git stash apply`:** Εφαρμόζει τη δουλειά, αλλά την κρατάει και στο stash (χρήσιμο αν θες να την εφαρμόσεις σε πολλά branches).

## 10. Επίλυση Διενέξεων (Merge Conflicts)

Όταν δύο προγραμματιστές (ή εσύ σε διαφορετικά branches) αλλάζουν τις *ίδιες ακριβώς γραμμές* σε ένα αρχείο, το Git δεν μπορεί να αποφασίσει ποια αλλαγή να κρατήσει. Αυτό δημιουργεί ένα **Merge Conflict**.

**Βήματα για ασφαλή επίλυση (π.χ. κατά τη διάρκεια ενός rebase):**
1. Το Git σταματάει τη διαδικασία και σε ειδοποιεί. Τρέξε `git status` για να δεις με κόκκινο χρώμα ποια αρχεία έχουν conflict.
2. Άνοιξε τα προβληματικά αρχεία στον Editor σου (π.χ. VS Code). Θα δεις τον κώδικα διαχωρισμένο κάπως έτσι:
   ```text
   <<<<<<< HEAD (Current Change - Ο κώδικας που υπήρχε ήδη)
   fmt.Println("Hello User")
   =======
   fmt.Println("Hello Admin")
   >>>>>>> feature-branch (Incoming Change - Ο κώδικας που προσπαθείς να βάλεις)
   ```
3. Επίλεξε ποιον κώδικα θέλεις να κρατήσεις, διαγράφοντας τα σύμβολα του Git (`<<<<<<<`, `=======`, `>>>>>>>`). Το VS Code έχει ενσωματωμένα κουμπιά "Accept Current", "Accept Incoming" ή "Accept Both" που το κάνουν αυτόματα.
4. Μόλις διορθώσεις όλα τα αρχεία και βεβαιωθείς ότι ο κώδικας είναι σωστός, τρέξε `git add <όνομα_αρχείου>` για να τα μαρκάρεις ως επιλυμένα (resolved).
5. Ολοκλήρωσε τη διαδικασία τρέχοντας `git rebase --continue` (αν έκανες rebase) ή `git commit` (αν έκανες merge). **Ποτέ μην τρέξεις `git add .` στα τυφλά κατά τη διάρκεια ενός conflict!**

## 11. Pull Requests & Code Reviews

Η ανάπτυξη σε ομάδες (ή ακόμα και μόνος σου για διατήρηση ποιότητας) απαιτεί μια διαδικασία ελέγχου πριν ο κώδικας μπει στο `main` branch.

1. **Άνοιγμα Pull Request (PR):** Μόλις ολοκληρώσεις και "καθαρίσεις" (rebase) το branch σου, το κάνεις push στο remote (π.χ. GitHub/Gitea) και ανοίγεις ένα PR προς το `main`.
   ```bash
   # Στείλε το τοπικό σου branch στο remote repository
   git push -u origin feat/add-login-page
   ```
2. **Code Review:** Εδώ μπαίνει σε εφαρμογή όλος ο οδηγός **Good Practices**! Ένας άλλος προγραμματιστής (ή εσύ σε δεύτερο χρόνο) διαβάζει τον κώδικα ελέγχοντας:
   - Τηρούνται οι αρχές SOLID και το Clean Code;
   - Είναι ο κώδικας "ιδιωματικός" (Effective Go);
   - Είναι πράσινα τα αυτόματα tests και οι έλεγχοι (π.χ. από το `check.sh`);
3. **Merge (Συγχώνευση):** Όταν εγκριθεί, το PR γίνεται merge. Προτίμησε τη μέθοδο **"Squash and Merge"** (αν δεν έκανες interactive rebase τοπικά), ώστε τα commits του branch σου να μπουν στο `main` ως ένα, καθαρό και περιεκτικό commit.

## 12. Tags και Εκδόσεις (Releases)

Όταν ένα project ή ένα σημαντικό ορόσημο (milestone) ολοκληρώνεται, είναι εξαιρετική πρακτική να βάζεις ένα Tag (ένα "στιγμιότυπο" του κώδικα στο χρόνο).

```bash
git tag v1.0.0 -m "Final submission for project"
git push origin v1.0.0
```

---

## 13. Προχωρημένα Εργαλεία Git (Advanced Tools)

Για τους Senior Developers, το Git προσφέρει εργαλεία "χειρουργικής" ακρίβειας για διάσωση κώδικα και debugging.

### Α. Git Reflog (Το Απόλυτο Δίχτυ Ασφαλείας)
Έκανες `git reset --hard` ή διέγραψες ένα branch κατά λάθος και έχασες τον κώδικά σου; Το Git σπάνια διαγράφει κάτι οριστικά. Το `reflog` καταγράφει **κάθε** κίνηση του HEAD, ακόμα και αυτές που δεν φαίνονται στο `git log`.
```bash
# Δες το ιστορικό όλων των ενεργειών σου
git reflog

# Βρες το hash της ενέργειας πριν κάνεις το λάθος (π.χ. HEAD@{2})
# Και επανέφερε τον κώδικά σου σε αυτό το σημείο!
git reset --hard HEAD@{2}
```

### Β. Git Cherry-Pick (Επιλεκτική Μεταφορά)
Έχεις ένα branch `feature-A` και ένα `feature-B`. Θέλεις ΜΟΝΟ ένα συγκεκριμένο commit (π.χ. ένα bugfix) από το `A` στο `B`, χωρίς να κάνεις merge όλο το branch.
```bash
# Πήγαινε στο branch που θέλεις να λάβει τον κώδικα
git checkout feature-B

# Εφάρμοσε μόνο το συγκεκριμένο commit
git cherry-pick <commit-hash>
```

### Γ. Git Bisect (Εντοπισμός Bugs με Δυαδική Αναζήτηση)
Κάποιο commit τις τελευταίες εβδομάδες εισήγαγε ένα bug, αλλά δεν ξέρεις ποιο; Το `bisect` κάνει δυαδική αναζήτηση στο ιστορικό, μειώνοντας τον χρόνο εντοπισμού στο μισό.
```bash
git bisect start
git bisect bad                 # Το τωρινό commit έχει το bug
git bisect good <commit-hash>  # Ένα παλιότερο commit που δούλευε σωστά

# Το Git θα κάνει checkout ένα commit στη μέση. 
# Ελέγχεις τον κώδικα (π.χ. τρέχεις τα tests σου) και του λες:
git bisect bad  # (αν το bug υπάρχει εδώ) ή git bisect good (αν δεν υπάρχει)

# Επαναλαμβάνεις μέχρι το Git να σου πει ποιο ακριβώς commit έσπασε τον κώδικα!
git bisect reset # Για να τερματίσεις τη διαδικασία και να γυρίσεις στο branch σου
```

## 14. Αυτοματοποίηση με Git Hooks

Μπορούμε να εξασφαλίσουμε ότι κανένας κακός κώδικας δεν θα γίνει ποτέ commit ή push! Το Git μας επιτρέπει να τρέχουμε αυτόματα scripts **πριν** ολοκληρωθεί μια ενέργεια (Git Hooks).

Μπορείς να συνδέσεις το `check.sh` με το Git μέσω ενός Pre-commit hook, ώστε ο έλεγχος να γίνεται αυτόματα πριν από κάθε commit. Για οδηγίες ρύθμισης, δες το έγγραφο **Automation and Scripts** και τον οδηγό για τα Conventional Commits.

---

**💡 Tip - Παράκαμψη των Hooks:**  
*Αν βρίσκεσαι σε ένα δικό σου branch και χρειάζεται να κάνεις ένα επείγον, πρόχειρο commit ή push (π.χ. Work In Progress) που ξέρεις ότι "σπάει" τους ελέγχους, μπορείς να παρακάμψεις τα `pre-commit` ή `pre-push` hooks προσθέτοντας το flag `--no-verify` (ή `-n`\):*

```bash
git commit -m "WIP: save progress" --no-verify
git push origin my-feature --no-verify
```

*(Χρησιμοποίησέ το με σύνεση και ποτέ στο `main` branch!)*
