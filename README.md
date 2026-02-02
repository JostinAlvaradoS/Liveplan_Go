# Liveplan Backend - Go

## ¿Qué es Liveplan?

Liveplan es una **plataforma de evaluación y análisis de proyectos de inversión** que permite a emprendedores y empresas:

- Crear planes de negocio detallados con proyecciones financieras
- Evaluar la viabilidad económica de sus proyectos
- Calcular indicadores clave (VAN, TIR, ROI)
- Realizar análisis de sensibilidad para entender riesgos
- Modelar diferentes escenarios de negocio

## Propósito del Backend

Este backend en Go actúa como el **motor de cálculos financieros** de Liveplan. Su responsabilidad principal es:

1. **Almacenar datos** del proyecto: inversiones, costos, ingresos, financiamiento
2. **Ejecutar cálculos complejos**: flujos de efectivo, balances, evaluaciones
3. **Mantener consistencia**: recalcular automáticamente cuando cambian datos
4. **Proveer información**: endpoints para consultar resultados

## Flujo Financiero Principal

El sistema calcula de forma **secuencial y dependiente**:

```
Datos Base (Productos, Precios, Costos, Inversiones)
    ↓
Presupuesto de Ventas (unidades × precio)
    ↓
Costos de Venta (costo variable × unidades)
    ↓
Gastos Operacionales (salarios, servicios, etc)
    ↓
Depreciaciones y Amortizaciones (desgaste de activos)
    ↓
Estado de Resultados (ingresos - gastos = utilidad)
    ↓
Flujo de Efectivo (cuando realmente entra/sale dinero)
    ↓
Balance General (posición financiera)
    ↓
Evaluación (VAN, TIR, indicadores)
```

## Módulos Principales

### 1. **Gestión de Datos de Proyecto**
Almacena la información base del plan de negocio:
- Productos/servicios a vender
- Precios y costos unitarios
- Inversiones iniciales
- Datos de préstamos y financiamiento
- Supuestos económicos (inflación, crecimiento)

### 2. **Cálculo de Ingresos**
- **Presupuesto de Venta**: Proyecta unidades vendidas por mes (crecimiento, estacionalidad)
- **Ventas en Dinero**: Aplica precios a las unidades
- **Políticas de Venta**: Define qué porcentaje es crédito vs contado (afecta flujo)

### 3. **Cálculo de Costos**
- **Costo de Materias Primas**: Costo variable directo
- **Costos de Producción**: Otros costos directos
- **Costos de Venta**: Comisiones, empaque, distribución
- **Gastos Operacionales**: Sueldos, arriendo, servicios
- **Depreciaciones y Amortizaciones**: Pérdida de valor de activos

### 4. **Financiamiento**
- **Datos de Préstamos**: Monto, tasa, plazo
- **Cálculo de Cuotas**: Amortización mensual con intereses
- **Políticas de Compra**: Define crédito vs contado en compras (afecta flujo)

### 5. **Estados Financieros**
- **Estado de Resultados**: Ingresos menos gastos = utilidad neta (accrual)
- **Flujo de Efectivo**: Movimiento real de dinero (cuando entra/sale)
- **Balance General**: Situación financiera del proyecto

### 6. **Evaluación del Proyecto**
- **VAN (Valor Actual Neto)**: Valor hoy de todos los flujos futuros
- **TIR (Tasa Interna Retorno)**: Rentabilidad del proyecto
- **Período de Recuperación**: Cuánto tarda en recuperarse la inversión

### 7. **Análisis de Sensibilidad**
Responde: **¿Qué pasa si los supuestos cambian?**
- Matriz 7×7 (49 escenarios)
- Varia volumen de ventas (-15% a +15%)
- Varia costos de producción (-15% a +15%)
- Calcula VAN para cada combinación
- Identifica cuál variable impacta más

## Arquitectura Técnica

### Stack Tecnológico
- **Lenguaje**: Go (compilado, rápido, seguro)
- **ORM**: GORM (abstracción de datos)
- **Base de Datos**: PostgreSQL (transacciones ACID)
- **API**: REST con net/http estándar

### Estructura de Código

```
internal/
├── models/          Definición de datos (estructuras Go)
├── controllers/     Lógica de negocio (cálculos)
├── handlers/        Ruteo HTTP (endpoints)
├── procedimientos/  Cálculos complejos
│   ├── recalcular.go        Cadena principal
│   ├── recalcular2.go       Análisis de sensibilidad
│   ├── calcular_prestamo.go Amortización
│   └── ...otros cálculos
└── db/              Conexión a BD
```

### Patrón de Recalculación

**Problema:** Si cambias un precio, ¿qué más hay que recalcular?
- Ventas en dinero
- Estado de resultados
- Flujo de efectivo
- Evaluación del proyecto

**Solución:** Cadena de dependencias

```
Recalcular(db, planID)
├─ CalcularPreciosYCostos()
├─ CalcularPresupuestoVenta()
├─ CalcularCostosVenta()
├─ CalcularGastosOperacion()
├─ CalcularDepreciaciones()
├─ CalcularEstadoResultados()
├─ CalcularFlujoEfectivo()
├─ CalcularBalanceGeneral()
└─ CalcularEvaluacionProyecto()
```

Cuando modificas algo, el flag `recalc: true` ejecuta toda esta cadena automáticamente.

## Ejemplo: Crear un Plan de Negocio

1. **Crear Plan**: Especificar duración (5 años), moneda, etc
2. **Definir Productos**: Qué venderá y a qué precio
3. **Especificar Costos**: Materias primas, operacionales
4. **Inversión Inicial**: Equipos, infraestructura
5. **Financiamiento**: Si necesita préstamo
6. **Proyecciones**: Crecimiento de ventas, inflación
7. **Políticas**: Qué % de ventas/compras son crédito

**Sistema calcula automáticamente:**
- Flujo de efectivo mes a mes por 5 años
- Si es viable (VAN > 0)
- Cuándo se recupera la inversión
- Sensibilidad a cambios en volumen/costos

## Características Clave

### ✅ Recalculación Automática
Cambias un dato → Todo se recalcula al instante

### ✅ Análisis de Sensibilidad
49 escenarios para entender riesgos y oportunidades

### ✅ Políticas Flexibles
Define crédito/contado mes a mes para venta y compra

### ✅ Depreciaciones
Soporta 2 tipos: depreciación contable y amortización de intangibles

### ✅ Préstamos Adaptativos
La cuota se recalcula si cambias monto, tasa o plazo

### ✅ Estados Financieros Reales
Diferencia entre accrual (contable) y efectivo (tesorería)

## Casos de Uso

### Para Emprendedores
"Quiero saber si mi negocio de café es viable con una inversión de $50k"

### Para Inversores
"¿Cuál es el retorno esperado de este proyecto en 5 años?"

### Para Bancos
"¿Puede este proyecto generar flujo para pagar un préstamo?"

### Para Analistas
"¿Cuál variable impacta más? ¿Precio o volumen?"

## Instalación y Uso

### Requisitos
- Go 1.20+
- PostgreSQL 12+
- Docker (opcional)

### Ejecutar
```bash
docker compose up -d      # Levanta PostgreSQL
go build
./liveplan_backend_go     # Servidor en :8080
```

### Variables de Entorno
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tu_contraseña
DB_NAME=liveplan
BACKEND_PORT=8080
```

## Changelog Reciente

### Febrero 2026
- ✅ Agregadas Políticas de Venta y Compra (mapeo crédito/contado por mes)
- ✅ Implementada lógica `recalc` para recalculación automática
- ✅ Corregida cuota de préstamo para adaptarse a cambios de parámetros
- ✅ Validada integración en flujo de efectivo

## Notas Técnicas

- Los cálculos son **determinísticos**: mismos datos = mismos resultados
- La cadena de recalculación es **secuencial** para mantener consistencia
- Los períodos son **años (1-5)** con desglose mensual
- Las tasas son **anuales** (se convierten a mensual internamente)
- Los montos están en **unidades monetarias** sin símbolo específico

---

**Liveplan Backend** es el corazón financiero de la plataforma. Procesa datos, ejecuta cálculos y proporciona análisis para que los emprendedores tomen decisiones informadas sobre sus proyectos.
