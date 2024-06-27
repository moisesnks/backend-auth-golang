# Backend de Autenticación de Usuarios (Golang)

Este repositorio contiene la implementación de un microservicio de autenticación de usuarios utilizando Firebase y otros servicios en la nube. El proyecto fue desarrollado por mí, Moisés Leiva, como parte de mi trabajo como Desarrollador Full-stack en la empresa ficticia UtemTX, creada por los estudiantes de TSSW20241S `(Taller de Ingeniería de Software 2024 1S)`, en el cual también cumplí mi rol de ayudante del ramo.

## Funcionalidades Implementadas

- **/login**: Permite a los usuarios iniciar sesión utilizando Firebase Auth y devuelve un token JWT válido.
- **/register**: Registra nuevos usuarios utilizando Firebase Auth y guarda información adicional en Firestore.
- **/verify**: Verifica códigos de verificación enviados por correo electrónico durante el registro.
- **/forgotPassword**: Permite a los usuarios solicitar recuperación de contraseña mediante correo electrónico.
- **/updateUser**: Permite a los usuarios actualizar su información personal.
- **/uploadPhoto**: Permite a los usuarios cargar y actualizar su foto de perfil.
- **/validateToken**: Valida tokens JWT y devuelve la información del usuario.

## Requisitos

- Implementación utilizando Firebase Auth para la autenticación.
- Almacenamiento de datos adicionales en Firestore.
- Almacenamiento de fotos de perfil en Firebase Storage.
- Envío de correos electrónicos utilizando SMTP para verificaciones y recuperación de contraseña.
- Documentación de los endpoints utilizando Swagger.
- Manejo de errores y respuestas consistentes.

## Estado Actual

- Implementación inicial completada con integración de Firebase Auth, Firestore, Firebase Storage y SMTP.
- Funcionalidades básicas como registro, verificación, actualización de usuario y carga de fotos de perfil están completamente funcionales.
- Documentación Swagger disponible localmente para consulta.

## Próximos Pasos

- Finalizar pruebas unitarias e integración para asegurar la fiabilidad y funcionalidad completa de cada endpoint.
