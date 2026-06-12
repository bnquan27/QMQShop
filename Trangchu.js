
const sliderImages = [
    "https://images.unsplash.com/photo-1603481588273-2f908a9a7a1b?q=80&w=1200&auto=format&fit=crop",
    "https://images.unsplash.com/photo-1547082299-de196ea013d6?q=80&w=1200&auto=format&fit=crop"
];
let currentSlideIndex = 0;
const imgElement = document.getElementById("slider-img");

function nextSlide() {
    currentSlideIndex = (currentSlideIndex + 1) % sliderImages.length;
    imgElement.src = sliderImages[currentSlideIndex];
}

function prevSlide() {
    currentSlideIndex = (currentSlideIndex - 1 + sliderImages.length) % sliderImages.length;
    imgElement.src = sliderImages[currentSlideIndex];
}

setInterval(nextSlide, 5000);


let cartCount = 0;
function addToCart() {
    cartCount++;
    document.querySelector(".cart-count").innerText = cartCount;
    alert("Đã thêm sản phẩm thành công vào giỏ hàng!");
}