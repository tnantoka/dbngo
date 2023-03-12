const go = new Go();

(async () => {
  const result = await WebAssembly.instantiateStreaming(fetch('main.wasm'), go.importObject);
  const instance = result.instance;
  go.run(instance);

  await Promise.all(examples.map(async(example) => {
    const response = await fetch(`examples/${example.name}`);
    example.content = await response.text();
  }));

  const select = document.querySelector('[data-js="select"]');
  const text = document.querySelector('[data-js="text"]');
  const run = document.querySelector('[data-js="run"]');
  const image = document.querySelector('[data-js="image"]');
  const errors = document.querySelector('[data-js="errors"]');

  examples.forEach((example) => {
    const option = document.createElement('option');
    option.value = example.name;
    option.textContent = example.name;
    select.appendChild(option);
  });

  select.addEventListener('change', () => {
    const example = examples.find((example) => example.name === select.value);
    text.value = example.content;
  });
  select.dispatchEvent(new Event('change'));

  run.addEventListener('click', async (e) => {
    e.preventDefault();

    image.src = '';
    errors.textContent = '';

    const generated = generate(text.value);
    console.log(generated);

    if (/data:/.test(generated)) {
      image.src = generated;
    } else {
      errors.textContent = generated;
    }
  });
})();
