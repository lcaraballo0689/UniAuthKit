📚 UniAuthKit – Guía Para Dummies (Completa)

> Objetivo: que cualquier persona —sin experiencia en Go— pueda levantar, probar y personalizar el microservicio UniAuthKit en < 30 min.




---

Índice rápido

1. Conceptos clave


2. Instalación Express


3. Configurar sin romper nada


4. Flujo de usuario


5. Comandos cURL


6. Cambiar nombres JSON


7. FAQs


8. Plantilla de proyecto real


9. Errores comunes


10. Despegue




---

0 🔰 Conceptos Clave


---

1 ⚡ Instalación Express

# Descarga el constructor
wget -O setup.sh https://raw.githubusercontent.com/lcaraballo0689/UniAuthKit/main/setup_login_microservice_full.sh
chmod +x setup.sh

# Ejecuta y genera estructura
./setup.sh

# Arranca servidor
authkit=$(pwd)/login-service
cd "$authkit"
go run ./cmd

Salida esperada:

Fiber v2 ▶ Listening on :3000


---

2 📝 Configurar sin romper nada (config/config.yaml)

<small>Piensa en este archivo como un panel de opciones.</small>

2.1 Cambiar puerto

server:
  port: "4000"

2.2 Cambiar motor BD a Postgres

database:
  driver: "postgres"
  dsn: "host=localhost user=dev password=dev dbname=login port=5432 sslmode=disable"

(createdb login si no existe.)

2.3 Activar/Desactivar módulos

endpoints:
  register: false  # oculta /register (403)
  oauth:    false  # desactiva OAuth Google
  mfa:      true   # activa MFA demo

2.4 Modo Súper‑Admin instantáneo

queries:
  permisos_usuario: ""   # vacío → role=superadmin, perms=["*"]


---

3 🔐 Flujo Básico de Usuario

[Cliente] → /register  → contraseña almacenada bcrypt
           → /login    → compara bcrypt
                  ↳ permisos_usuario?  Sí → lista, No → superadmin/*
           ← JSON con access_token (30 min) + refresh_token (7 días)


---

4 🎮 Comandos cURL (copiar / pegar)

# Registrar
curl -X POST :3000/register -H 'Content-Type: application/json' \
  -d '{"username":"ana","password":"1234","email":"ana@mail"}'

# Login → guarda tokens
TOKENS=$(curl -s -X POST :3000/login -H 'Content-Type: application/json' \
  -d '{"username":"ana","password":"1234"}')
ACCESS=$(echo $TOKENS | jq -r .access_token)
REFRESH=$(echo $TOKENS | jq -r .refresh_token)

echo "Access:  $ACCESS"
echo "Refresh: $REFRESH"

# Revocar refresh
curl -X POST :3000/token/revoke -H 'Content-Type: application/json' \
  -d "{\"token\":\"$REFRESH\"}"

# MFA Setup
curl -X POST :3000/mfa/setup -H 'Content-Type: application/json' \
  -d '{"userID":1}'

# MFA Verify (usa el secreto devuelto como código)
curl -X POST :3000/mfa/verify -H 'Content-Type: application/json' \
  -d '{"userID":1,"code":"<SECRET>"}'


---

5 🛠️ Cambiar nombres del JSON

response_template:
  status_field: "estado"
  code_field:   "codigo"
  message_field:"mensaje"

Ahora la respuesta será:

{"estado":"success","codigo":200,"mensaje":"👍"}


---

6 🧑‍🎓 FAQs

“Falla conexión BD” → revisa driver/dsn, usa sqlite de prueba.

“No veo /register” → endpoints.register = false.

“invalid_credentials” → contraseña incorrecta (cuidado espacios).

“open config.yaml” → lanza desde raíz o usa flag -config.



---

7 🗂️ Plantilla de Proyecto real

1. Cambia secret JWT antes de producción ✔️


2. Sustituye OTP demo por TOTP real (pquerna/otp).


3. Ajusta SQL permisos.


4. Añade tabla intentos (anti‑brute‑force).


5. Empaqueta con Docker Compose.




---

8 🚧 Errores comunes


---

9 🏁 ¡Listo para despegar!

Modifica YAML 🠚 cambia comportamiento sin recompilar.

Usa superadmin solo en desarrollo 🛑.

Avanza módulo por módulo (internal/*).


Si algo falla: vuelve a sqlite + desactiva extras, prueba de nuevo.

¡Hackea feliz! 🎉

