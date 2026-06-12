function switchTab(tabName) {

    const tabs = document.querySelectorAll('.tab-btn');
    const forms = document.querySelectorAll('.auth-form');

   
    tabs.forEach(tab => tab.classList.remove('active'));
    forms.forEach(form => form.classList.remove('active'));

    if (tabName === 'login') {
        tabs[0].classList.add('active');
        document.getElementById('login-form').classList.add('active');
    } else if (tabName === 'register') {
        tabs[1].classList.add('active');
        document.getElementById('register-form').classList.add('active');
    }
}