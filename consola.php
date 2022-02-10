<?php
/*conexion a la base de dato*/
function conectar(){
	$servidor = "127.0.0.1";
	$usuario = "root";
	$password = "123456";
	$database = "chat";
	$conexion = new mysqli($servidor, $usuario, $password, $database);
	return $conexion;
}
/*funcion para ejecutar*/
function ejecutar($sql){
	$conexion = conectar();
	mysqli_query($conexion, $sql);
}
/*funcion para consultar*/
function consultar($sql,$cols_num){
	$conexion = conectar();
	$query = $conexion->query($sql);
	$matriz = array();
	$f = 0;
	while($celda = $query->fetch_assoc()){
		$keys = array_keys($celda);
		for($c=0;$c<$cols_num;$c++){$matriz[$f][$c]=$celda[$keys[$c]];}
		$f++;
	}
	return $matriz;
}
/*funcion de retorno de datos AJAX*/
function AJAX($nombre, $mensaje){
	if(($nombre!="")&&($mensaje!="")){
		ejecutar("INSERT INTO mensajeria(usuarios,mensajes) VALUES('".$nombre."','".$mensaje."');");
	}
	$chat = consultar("SELECT concat(idmensajeria, ';', usuarios, ';', mensajes)FROM mensajeria ORDER BY idmensajeria DESC LIMIT 5 ",1);
	$i = 0;
	$caracteres = "";
	foreach ($chat as $dato){
		if($i == 0){
			$caracteres = preg_replace("/\r|\n/", "-n", $dato[0]);
		}
		else{
			$caracteres = preg_replace("/\r|\n/", "-n", $dato[0])."\n". $caracteres;
		}
		$i = $i+1;
	}
	header("Content-Type: text/plain");
	echo($caracteres);
}
/*solo si recive variable nombre y mensaje sabemos es el AJAX*/
if(isset($_REQUEST["nombre"])&&isset($_REQUEST["mensaje"])){
	$nombre = $_REQUEST["nombre"];
	$mensaje = $_REQUEST["mensaje"];
	AJAX($nombre, $mensaje);
}
/*sin intenta ingresar sin autorizacion*/
else{
	echo("Solo Personal Autorizado");
}
?>
