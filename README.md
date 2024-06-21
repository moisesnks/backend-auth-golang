### Ticket: Implementación Completa del Microservicio de Autenticación de Usuarios

## Descripción
Se requiere la implementación de un microservicio de autenticación que soporte las siguientes funcionalidades básicas y adicionales:

1. **/login** - Inicio de sesión de usuarios.
2. **/register** - Registro de usuarios.
3. **/verify** - Verificación de códigos de verificación enviados por correo.
4. **/forgotPassword** - Recuperación de contraseña.
5. **/updateUser** - Actualización de información del usuario.
6. **/uploadPhoto** - Carga de fotos de perfil del usuario.
7. **/validateToken** - Validación de tokens de autenticación y devolución de la información del JWT.

## Requisitos
- [x] **/login** - Handler para el inicio de sesión de usuarios:
   - Validar las credenciales del usuario utilizando Firebase Auth.
   - Devolver un token JWT al usuario autenticado.

- [x] **/register** - Handler para el registro de usuarios:
   - Registrar nuevos usuarios utilizando Firebase Auth.
   - Guardar información adicional del usuario en Firestore.
   - Asignar el rol de 'guest' al nuevo usuario.
   - Enviar un correo de verificación con un código utilizando SMTP.

- [x] **/verify** - Handler para la verificación de códigos de verificación:
   - Verificar el código recibido en la solicitud con el código almacenado en Firestore.
   - Cambiar el rol del usuario a 'member' si el código es correcto.

- [ ] **/forgotPassword** - Handler para la recuperación de contraseña:
   - Generar un token de recuperación de contraseña utilizando Firebase Auth.
   - Enviar el token al correo electrónico del usuario utilizando SMTP.

- [ ] **/updateUser** - Handler para la actualización de información del usuario:
   - Permitir a los usuarios autenticados actualizar su información personal en Firestore.
   - Asegurar que solo el usuario autenticado pueda actualizar su propia información.

- [ ] **/uploadPhoto** - Handler para la carga de fotos de perfil del usuario:
   - Permitir a los usuarios autenticados cargar y actualizar su foto de perfil en Firebase Storage.
   - Asegurar que solo el usuario autenticado pueda cargar y actualizar su propia foto de perfil.

- [ ] **/validateToken** - Handler para la validación de tokens de autenticación:
   - Validar tokens JWT utilizando Firebase Auth.
   - Devolver la información contenida en el token JWT.

## Implementación
- [x] Utilizar Firebase Auth para la autenticación y gestión de usuarios.
- [x] Utilizar Firestore para almacenar información adicional de los usuarios.
- [x] Utilizar Firebase Storage para almacenar fotos de perfil de los usuarios.
- [x] Utilizar SMTP para el envío de correos electrónicos.
- [x] Documentar cada endpoint utilizando Swagger.
- [ ] Manejar adecuadamente los errores y proporcionar respuestas claras y consistentes.
- [ ] Realizar pruebas unitarias y de integración para asegurar la funcionalidad y fiabilidad de cada endpoint.

### Actualizaciones 
- He completado la integración inicial de Firebase Auth, Firestore, Firebase Storage y SMTP para todas las funcionalidades de usuario.
- Se ha implementado correctamente el envío de correo de verificación durante el registro en /register.
- La verificación de códigos de verificación en /verify está completamente funcional.
- Estoy trabajando en la implementación de la recuperación de contraseña, actualización de información de usuario y carga de fotos de perfil.
- La documentación Swagger está disponible para consulta local en [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html).

### Observación 
_Con lo que llevo, el equipo de infraestructura ya podría deployar este microservicio. Durante el sprint, continuaré con el CI/CD para completarlo. Esto permitirá que el equipo de frontend y QA participen en las integraciones en el cliente y en las pruebas de calidad. Por favor, contáctense conmigo en mleiva@utem.cl para enviarles las variables de entorno necesarias en caso de deploy en infraestructura, y para conocer más sobre el código y su integración en el frontend._

**Asignado a:** 
[moisesnks](github.com/moisesnks)
**Prioridad:**
Alta
