import fs from 'fs';

// Variable definitions follow this format:
// { name: [light, dark, black] }

// Create a pseudo-3D border using subpixel box-shadow with no blur.
const pseudo3DBorder = (
  x: number,
  y: number,
  thickness: number,
  color: string,
  alpha: string,
) => `${x}px ${y}px 0 ${thickness}px #${color}${alpha}`;

// Convenience function for light-mode (dark shadow) pseudo-3D borders.
const lightP3D = (
  alpha = '40',
  y = 0.4,
  x = 0,
  thickness = 0.8,
  color = '000000',
) => pseudo3DBorder(x, y, thickness, color, alpha);

// Convenience function for dark-mode (light highlight) pseudo-3D borders.
const darkP3D = (
  alpha = '30',
  y = -0.33,
  x = 0,
  thickness = 1,
  color = 'ffffff',
) => pseudo3DBorder(x, y, thickness, color, alpha);

// Compensate for Firefox's incorrect subpixel rendering with a slight 1px thick border.
const ffRenderFix = (alpha: string, color: string) =>
  `0 0 0 1px #${color}${alpha}`;
const lightFFRF = (alpha = '12') => ffRenderFix(alpha, '000000');
const darkFFRF = (alpha = '15') => ffRenderFix(alpha, 'ffffff');

// Add " inset".
const inset = (s: string) => s + ' inset';

// Join with commas.
const list = (...s: string[]) => s.join(', ');

const themeVariables = {
  'primary-color': ['#000000', '#ffffff', '#ffffff'],
  'strong-color': ['#000000', '#ffffff', '#ffffff'],
  'grey-1-color': ['#111111', '#eeeeee', '#eeeeee'],
  'grey-2-color': ['#222222', '#dddddd', '#dddddd'],
  'grey-3-color': ['#333333', '#cccccc', '#cccccc'],
  'grey-4-color': ['#444444', '#bbbbbb', '#bbbbbb'],
  'grey-5-color': ['#555555', '#aaaaaa', '#aaaaaa'],
  'grey-6-color': ['#666666', '#999999', '#999999'],
  'grey-7-color': ['#777777', '#888888', '#888888'],
  'grey-8-color': ['#888888', '#777777', '#777777'],
  'grey-9-color': ['#999999', '#666666', '#666666'],
  'grey-a-color': ['#aaaaaa', '#555555', '#555555'],
  'grey-b-color': ['#bbbbbb', '#444444', '#444444'],
  'grey-c-color': ['#cccccc', '#333333', '#333333'],
  'grey-d-color': ['#dddddd', '#222222', '#222222'],
  'grey-e-color': ['#eeeeee', '#111111', '#111111'],
  'grey-e7-color': ['#e7e7e7', '#0c0c0c', '#0c0c0c'],
  'grey-f9-color': ['#f9f9f9', '#070707', '#070707'],
  'layout-bg-color': ['#cccccc', '#222222', '#222222'],
  'secondary-bg-color': ['#f5f5f5', '#050505', '#050505'],
  'primary-bg-color': ['#ffffff', '#000000', '#000000'],

  // Contextual color
  'error-color': ['#d00000', '#ff1111', '#ff1111'],
  'error-color-faded': ['#ff1111', '#d00000', '#d00000'],
  'link-color': ['#0066ee', '#22aaff', '#22aaff'],
  'online-green-color': ['#00c220', '#00c220', '#00c220'],

  // Danger buttons
  'danger-button-color': ['#ffffff', '#ffffff', '#ffffff'],
  'danger-button-background': [
    'linear-gradient(#ff0000, #cc0000)',
    'linear-gradient(#dd0000, #aa0000)',
    'linear-gradient(#dd0000, #aa0000)',
  ],
  'danger-button-hover-background': [
    'linear-gradient(#dd0000, #aa0000)',
    'linear-gradient(#bb0000, #880000)',
    'linear-gradient(#bb0000, #880000)',
  ],
  'danger-button-active-background': [
    'linear-gradient(#cc0000, #880000)',
    'linear-gradient(#aa0000, #660000)',
    'linear-gradient(#aa0000, #660000)',
  ],
  'danger-button-shadow': [
    list('0 4px 8px -3px #00000026', '0 0.4px 0 0.8px #880000ff', lightFFRF()),
    list('0 -0.33px 0 1px #ff555533', darkFFRF()),
    list('0 -0.33px 0 1px #ff555533', darkFFRF()),
  ],
  'danger-button-hover-shadow': [
    list('0 6px 8px -4px #00000033', '0 0.4px 0 0.8px #660000ff', lightFFRF()),
    list('0 -0.33px 0 1px #ff555555', darkFFRF()),
    list('0 -0.33px 0 1px #ff555555', darkFFRF()),
  ],
  'danger-button-active-shadow': [
    list('0 4px 6px -3px #00000026', '0 0.4px 0 0.8px #660000ff', lightFFRF()),
    list('0 -0.33px 0 1px #ff555544', darkFFRF()),
    list('0 -0.33px 0 1px #ff555544', darkFFRF()),
  ],

  // Nav
  'nav-bg-color': ['#e7e7e7', '#050505', '#050505'],
  'nav-button-active-bg-color': ['#d5d5d5', '#111111', '#111111'],
  'nav-button-active-hover-bg-color': ['#cccccc', '#1c1c1c', '#1c1c1c'],

  // Code
  'code-bg-color': ['#44444411', '#aaaaaa11', '#aaaaaa11'],
  'code-shadow': [
    list('0 6px 9px -4px #00000033', '0 0.4px 0 0.8px #0000001a', lightFFRF()),
    list('0 0.33px 0 1px #ffffff26', darkFFRF()),
    list('0 0.33px 0 1px #ffffff26', darkFFRF()),
  ],

  // Process details
  'sparkline-stroke': ['#ee0000', '#ee0000', '#ee0000'],
  'sparkline-fill': ['#ee000044', '#ee000044', '#ee000044'],

  // Spinners
  'spinner-grey': ['#777777ff', '#777777ff', '#777777ff'],
  'spinner-grey-light': ['#77777733', '#77777733', '#77777733'],

  // Checkbox
  'checkbox-color': ['#777777ff', '#999999ff', '#999999ff'],
  'checkbox-bg-color': ['#77777700', '#99999900', '#99999900'],
  'checkbox-border-color': ['#77777777', '#99999977', '#99999977'],
  'checkbox-hover-color': ['#444444', '#777777', '#777777'],
  'checkbox-hover-bg-color': ['#77777711', '#99999922', '#99999922'],
  'checkbox-focus-bg-color': ['#77777722', '#99999944', '#99999944'],
  'checkbox-active-hover-bg-color': ['#55ccff22', '#19baff22', '#19baff22'],
  'checkbox-active-focus-bg-color': ['#55ccff44', '#19baff44', '#19baff44'],
  'checkbox-active-color': ['#008ac5', '#19baff', '#19baff'],
  'checkbox-active-border-color': ['#008ac577', '#19baff77', '#19baff77'],

  // Buttons
  'icon-button-bg-color': ['#00000000', '#eeeeee00', '#eeeeee00'],
  'icon-button-hover-bg-color': ['#00000010', '#eeeeee18', '#eeeeee18'],
  'icon-button-focus-bg-color': ['#00000018', '#eeeeee33', '#eeeeee33'],
  'button-background': [
    'linear-gradient(#fff, #f5f5f5)',
    'linear-gradient(#141414, #040404)',
    'linear-gradient(#141414, #040404)',
  ],
  'button-hover-background': [
    'linear-gradient(#fafafa, #e7e7e7)',
    'linear-gradient(#242424, #111111)',
    'linear-gradient(#242424, #111111)',
  ],
  'button-active-background': [
    'linear-gradient(#f7f7f7, #e0e0e0)',
    'linear-gradient(#111111, #000000)',
    'linear-gradient(#111111, #000000)',
  ],
  'button-inset-background': [
    'linear-gradient(#00000022, #00000011)',
    'linear-gradient(#222222, #333333)',
    'linear-gradient(#222222, #333333)',
  ],
  'button-shadow': [
    list('0 4px 8px -3px #00000022', lightP3D('40'), lightFFRF()),
    list(darkP3D('30'), darkFFRF()),
    list(darkP3D('30'), darkFFRF()),
  ],
  'button-hover-shadow': [
    list('0 6px 8px -4px #00000030', lightP3D('59'), lightFFRF()),
    list(darkP3D('50'), darkFFRF()),
    list(darkP3D('50'), darkFFRF()),
  ],
  'button-active-shadow': [
    list('0 4px 6px -3px #00000024', lightP3D('73'), lightFFRF()),
    list(darkP3D('40'), darkFFRF()),
    list(darkP3D('40'), darkFFRF()),
  ],
  'button-inset-shadow': [
    list(
      inset('0 4px 8px -3px #00000022'),
      inset(lightP3D('40')),
      inset(lightFFRF()),
    ),
    list(inset(darkP3D('40')), inset(darkFFRF())),
    list(inset(darkP3D('40')), inset(darkFFRF())),
  ],

  // Shadows
  'heavy-3d-box-shadow': [
    list('0 8px 12px -6px #0000004d', lightP3D('30', 0.5, 0, 1), lightFFRF()),
    list(darkP3D('23'), '0 8px 12px -6px #0000004d', darkFFRF()),
    list(darkP3D('23'), '0 8px 12px -6px #0000004d', darkFFRF()),
  ],
  'text-input-shadow': [
    list(
      inset('0 6px 9px -4px #0000001a'),
      inset(lightP3D('1a')),
      inset(lightFFRF()),
    ),
    list(darkP3D('26'), inset('0 6px 9px -4px #0000001a'), inset(darkFFRF())),
    list(darkP3D('26'), inset('0 6px 9px -4px #0000001a'), inset(darkFFRF())),
  ],
  'text-input-shadow-focus': [
    list(
      '0 0px 0 1px #0066ee',
      inset('0 6px 9px -4px #0000001a'),
      inset(lightP3D('1a')),
      inset(lightFFRF()),
    ),
    list(
      '0 0px 0 1px #22aaff',
      inset('0 6px 9px -4px #0000001a'),
      inset(lightP3D('1a')),
      inset(darkFFRF()),
    ),
    list(
      '0 0px 0 1px #22aaff',
      inset('0 6px 9px -4px #0000001a'),
      inset(lightP3D('1a')),
      inset(darkFFRF()),
    ),
  ],
  'shadow-focus': [
    '0 0px 0 1px #0066ee',
    '0 0px 0 1px #22aaff',
    '0 0px 0 1px #22aaff',
  ],
  'dropdown-shadow': [
    list('0 12px 12px -4px #00000033', lightP3D('52'), lightFFRF()),
    list(darkP3D('55'), '0 12px 12px -4px #00000033', darkFFRF()),
    list(darkP3D('55'), '0 12px 12px -4px #00000033', darkFFRF()),
  ],
  'card-shadow': [
    list(
      '0 4px 8px -3px #00000022',
      '0.2px 0.3px 0 0.7px #00000030',
      lightFFRF(),
    ),
    list(
      '0.25px -0.25px 0 0.75px #ffffff30',
      '0 4px 8px -3px #00000022',
      darkFFRF(),
    ),
    list(
      '0.25px -0.25px 0 0.75px #ffffff30',
      '0 4px 8px -3px #00000022',
      darkFFRF(),
    ),
  ],
  'card-hover-shadow': [
    list(
      '0 6px 8px -4px #00000044',
      '0.2px 0.25px 0 0.85px #00000050',
      lightFFRF(),
    ),
    list(
      '0.25px -0.25px 0 1px #ffffff50',
      '0 6px 8px -4px #00000044',
      darkFFRF(),
    ),
    list(
      '0.25px -0.25px 0 1px #ffffff50',
      '0 6px 8px -4px #00000044',
      darkFFRF(),
    ),
  ],
};

const childVariables = {
  'log-color': [
    'var(--light-log-color)',
    'var(--dark-log-color)',
    'var(--dark-log-color)',
  ],
  'log-bg-color': [
    'var(--light-log-bg-color)',
    'var(--dark-log-bg-color)',
    'var(--dark-log-bg-color)',
  ],
  'log-hover-color': [
    'var(--light-log-hover-color)',
    'var(--dark-log-hover-color)',
    'var(--dark-log-hover-color)',
  ],
  'log-bg-hover-color': [
    'var(--light-log-bg-hover-color)',
    'var(--dark-log-bg-hover-color)',
    'var(--dark-log-bg-hover-color)',
  ],
};

// Format a list of root level theme variables.
const themeDefinition = (theme: number) =>
  Object.entries(themeVariables)
    .map((entry) => {
      return `  --${entry[0]}: ${entry[1][theme]};`;
    })
    .join('\n');

// Format a list of child level theme variables.
const childDefinition = (theme: number) =>
  Object.entries(childVariables)
    .map((entry) => {
      return `  --${entry[0]}: ${entry[1][theme]};`;
    })
    .join('\n');

// Format a block of theme variables for root and child selectors.
const themeBlock = (selector: string, theme: number) => `body.${selector} {
${themeDefinition(theme)}
}
body.${selector} * {
${childDefinition(theme)}
}`;

// Format the entire theme file using the above helpers.
const out = `/* Generated file. DO NOT EDIT. */
${themeBlock('auto', 0)}
@media (prefers-color-scheme: dark) {
${themeBlock('auto', 2)}}
${themeBlock('light', 0)}
${themeBlock('dark', 1)}
${themeBlock('black', 2)}
`;

// Save generated output.
fs.writeFile('./public/theme-generated.css', out, function (err) {
  if (err) throw err;
  console.log('Generated theme file.');
});
