# Implementacija konkurentnog chat servera u programskom jeziku Go

Istraživački projekat iz predmeta **Dizajn programskih jezika**  
Autor: Jelisaveta Gavrilović 1028/2024  
Datum: October 2025

---

## Opis projekta

Ovaj projekat predstavlja implementaciju konkurentnog chat servera u programskom jeziku Go. Sistem omogućava korisnicima razmenu poruka u realnom vremenu preko TCP protokola. Projekat demonstrira korišćenje gorutina i kanala za konkurentno izvršavanje, kao i mutex-a za zaštitu deljenih resursa. Terminalski interfejs omogućava dinamički prikaz poruka i korisničkog unosa, dok podrška za komande `/quit`, `/users` i `/help` olakšava interakciju. Privatne poruke se šalju u obliku `@username poruka`.

## Funkcionalnosti

- Terminalska komunikacija između korisnika  
- Slanje broadcast i privatnih poruka  
- Podrška za osnovne komande: `/quit`, `/users`, `/help`  
- Konkurentna obrada poruka za više korisnika istovremeno  
- Terminalski UI pomoću biblioteka `tview` i `tcell`  

## Instalacija

### 1. Instalacija Go okruženja

Preuzmite i instalirajte Go sa zvanične stranice: [https://go.dev/dl/](https://go.dev/dl/). Alternativno, možete koristiti sledeće komande u terminalu:

**Za Linux:**
```bash
sudo apt install golang-go

```

Za macOS:
```bash
brew install go
```

Provera instalirane verzije:
```bash
go version
```


### 2. Instalacija zavisnosti

Zavisnosti su definisane u go.mod (u `server` i `client` direktorijumu) i go.sum (u `client` direktorijumu). Go automatski preuzima potrebne biblioteke prilikom prvog pokretanja. Ukoliko želite, možete ih preuzeti ručno:
```bash
go get github.com/rivo/tview
go get github.com/gdamore/tcell/v2
```

Za osvežavanje i čišćenje zavisnosti, preporučuje se korišćenje:
```bash
go mod tidy
```

Ova komanda:  
- Preuzima sve nedostajuće zavisnosti  
- Briše nepotrebne zavisnosti  
- Održava go.mod i go.sum u konzistentnom stanju


## Preuzimanje projekta:

Projekat možete preuzeti pokretanjem sledeće komande:
```bash
git clone https://github.com/jelisavetagavrilovic/mini-chat-server
cd mini-chat-server
```

## Pokretanje

Server mora biti pokrenut pre klijenata.

Pokretanje servera:
```bash
go run server/main.go
```
Pokretanje klijenta:
```bash
go run client/main.go
```

Nakon pokretanja, klijent bira korisničko ime i može komunicirati sa ostalim korisnicima u realnom vremenu.



---







