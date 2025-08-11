function $(id){ return document.getElementById(id) }

async function fetchOrder(){
  const id = $('orderId').value.trim()
  if(!id){ alert('Введите order_uid'); return }
  $('result').textContent = 'Загрузка...'
  try {
    const res = await fetch(`/order/${encodeURIComponent(id)}`)
    const ct = res.headers.get('content-type') || ''
    if (ct.includes('application/json')) {
      const data = await res.json()
      $('result').textContent = JSON.stringify(data, null, 2)
    } else {
      const text = await res.text()
      $('result').textContent = text
    }
  } catch (e) {
    $('result').textContent = 'Ошибка: ' + e
  }
}

async function createOrder(){
  const raw = $('orderJson').value
  $('createResult').textContent = 'Отправка...'
  try{
    const res = await fetch('/order', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: raw })
    const ct = res.headers.get('content-type') || ''
    if (ct.includes('application/json')) {
      const data = await res.json()
      $('createResult').textContent = JSON.stringify(data, null, 2)
    } else {
      const text = await res.text()
      $('createResult').textContent = text
    }
  }catch(e){
    $('createResult').textContent = 'Ошибка: ' + e
  }
}

window.addEventListener('DOMContentLoaded', () => {
  $('btnFetch').addEventListener('click', fetchOrder)
  $('btnCreate').addEventListener('click', createOrder)
})


