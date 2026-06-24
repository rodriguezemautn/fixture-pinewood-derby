# Fixture para carreras 
 Fixture para carreras de pinewood derby cars
---

## Objetivo
Controlador de fixture de carrera de autitos derby de exploradores del Rey del Destacamento 15 de la iglesia Betel

## Alcance
Este sistema gestionara las carreras de cada categoria de autos de madera de pino con carreras de hasta 4 autos por instacia. Las categorias se dividiran por edad, participando una lista aun no determinada de autos. 

## Requerimientos Funcionales
1. El sistema debera poder registar las categorias de la fecha. 
2. El sistema debera poder registrar uno o varios autos por categoria. Con su foto, numero de auto, nombre de creador y edad. 
3. El sistema organizara las catedorias en carrera de maximo 4 autos por carrera. 
4. El sistema registrara el orden de llegada de cada auto en cada carrera.
5. El sistema schedulizara un fixture con una secuencia de carreras para obetener los 4 mejores competidores para la carrera final en cada categoria buscando la secuencia de carreras mas optimas y equitativa para todos los competidores.
6. El sistema registrara todo el fixture y lo mostrara en forma grafica y atractiva. 
7. El sistema anunciara de forma emocionante y visualmente desafiante el podio de cada carrera
8. El sistema anunciara en una ventana especial para el podio de final de cada carrera los ganadores.

## Requerimiento no funcionales

1. Mobile first 
2. Se proyectara en un proyecto con baja resolucion y en un monitor 
3. El sistema debera ser rapido de ejecucion 
4. El sistema podra ser accesible en modo visualizacion de manera online
5. El sistema debera persistir todo avance y recordar su estado aunque se interrumpa la conexion, electricidad o se interrumpa la ejecucion de forma abrupta.
6. El sistema debera poder ejecutarse como administrador de forma offline tambien.
7. El sistema debera ser portable 


## Reglas de Negocio

1. El sistema tendra un acceso de administrador quien gestionara las categorias y carreras. 
2. El sistema podra accederce como visualizador solo de forma online.


