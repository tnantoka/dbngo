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
  const imagePNG = document.querySelector('[data-js="image-png"]');
  const imageGIF = document.querySelector('[data-js="image-gif"]');
  const withGIF = document.querySelector('[data-js="with-gif"]');
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

  run.addEventListener('click', async (e) => {
    e.preventDefault();

    imagePNG.src = '';
    imageGIF.src = '';
    errors.textContent = '';

    const generatedPNG = generatePNG(text.value);
    if (/data:/.test(generatedPNG)) {
      imagePNG.src = generatedPNG;
      if (withGIF.checked) {
        imageGIF.src = generateGIF(text.value);
      }
    } else {
      errors.textContent = generatedPNG;
    }
  });

  select.value = 'mitpress/194-2.dbn';
  select.dispatchEvent(new Event('change'));
  run.dispatchEvent(new Event('click'));
})();
