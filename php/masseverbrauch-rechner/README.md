# Masseverbrauch-Rechner

Ein Tool zur Berechnung des Masseverbrauchs in der Pressabteilung einer Fliesenproduktion.

## Funktionsweise

- **Bandanlage**: Liefert Hauptmasse (z. B. 1,5 kg/s), gemessen mit Stoppuhr – Eingabe in Minuten erforderlich
- **Stellplätze**: Liefern Flakes; Erfassung des Gewichts vor und nach der Laufzeitmessung der Hauptmasse
- **Berechnung**: Bei Button-Klick wird der prozentuale Verbrauch jeder Flakes-Masse ermittelt (maximal 5 verschiedene Flakes-Massen, Gesamtsumme = 100 %)

## Technologie-Stack

- **Backend**: PHP mit Server-Side-Rendering und HTMX
- **Frontend**: HTML/CSS (ohne Frameworks)
- **JavaScript**: Minimaler Einsatz
- **API**: Golang-API zur Verarbeitung der Eingaben und Rückgabe der Ergebnisse (JSON)
