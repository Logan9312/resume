use axum::{
    http::{header, Method},
    response::{IntoResponse, Response},
    routing::get,
    Router,
};
use image::ImageFormat;
use pdfium_render::prelude::*;
use std::{env, io::Cursor, sync::OnceLock};
use tower_http::cors::{Any, CorsLayer};

static PDF_DATA: OnceLock<Vec<u8>> = OnceLock::new();
static PNG_DATA: OnceLock<Vec<u8>> = OnceLock::new();

fn init_data() {
    let pdf_bytes = std::fs::read("./resume.pdf").expect("Failed to read resume.pdf");
    PDF_DATA.set(pdf_bytes.clone()).unwrap();

    // Try to bind to PDFium (may not be available locally)
    let pdfium_result = Pdfium::bind_to_library("./libpdfium.so")
        .or_else(|_| Pdfium::bind_to_library("/app/libpdfium.so"))
        .or_else(|_| Pdfium::bind_to_library("./pdfium.dll"))
        .or_else(|_| Pdfium::bind_to_system_library());

    let pdfium = match pdfium_result {
        Ok(bindings) => Pdfium::new(bindings),
        Err(e) => {
            println!("PDFium not available, PNG endpoint disabled: {e}");
            return;
        }
    };

    let png_bytes = {
        let document = pdfium
            .load_pdf_from_byte_slice(&pdf_bytes, None)
            .expect("Failed to load PDF");

        let page = document.pages().get(0).expect("No pages in PDF");

        let bitmap = page
            .render_with_config(&PdfRenderConfig::new().set_target_width(1200))
            .expect("Failed to render page");

        let image = bitmap.as_image();
        let mut buf = Cursor::new(Vec::new());
        image
            .write_to(&mut buf, ImageFormat::Png)
            .expect("Failed to encode PNG");
        buf.into_inner()
    };

    PNG_DATA.set(png_bytes).unwrap();
    println!("PNG pre-generated successfully");
}

async fn serve_pdf() -> Response {
    (
        [(header::CONTENT_TYPE, "application/pdf")],
        PDF_DATA.get().unwrap().as_slice(),
    )
        .into_response()
}

async fn serve_png() -> Response {
    match PNG_DATA.get() {
        Some(data) => ([(header::CONTENT_TYPE, "image/png")], data.as_slice()).into_response(),
        None => (
            axum::http::StatusCode::SERVICE_UNAVAILABLE,
            "PNG not available - PDFium not loaded",
        )
            .into_response(),
    }
}

#[tokio::main]
async fn main() {
    init_data();

    let cors = CorsLayer::new()
        .allow_origin(Any)
        .allow_methods([Method::GET, Method::OPTIONS])
        .allow_headers([header::ORIGIN, header::CONTENT_TYPE, header::ACCEPT]);

    let app = Router::new()
        .route("/", get(serve_pdf))
        .route("/png", get(serve_png))
        .layer(cors);

    let port = env::var("PORT").unwrap_or_else(|_| "3000".to_string());
    let addr = format!("0.0.0.0:{}", port);

    println!("Listening on {}", addr);

    let listener = tokio::net::TcpListener::bind(&addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
