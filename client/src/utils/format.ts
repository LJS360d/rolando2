export function formatTime(seconds: number): string {
  const units = [
    { label: 'y', divisor: 60 * 60 * 24 * 30 * 12 },
    { label: 'm', divisor: 60 * 60 * 24 * 30 },
    { label: 'd', divisor: 60 * 60 * 24 },
    { label: 'h', divisor: 60 * 60 },
    { label: 'm', divisor: 60 },
    { label: 's', divisor: 1 },
  ];

  let output = '';
  for (const unit of units) {
    const value = Math.floor(seconds / unit.divisor);
    if (value > 0 || output !== '') {
      output += `${value.toString().padStart(2, '0')}${unit.label} `;
      seconds %= unit.divisor;
    }
  }

  return output.trim();
}

export function formatBytes(bytes: number | bigint): string {
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];

  if (bytes === 0) {
    return '0 Byte';
  }

  const i = Math.floor(Math.log2(Number(bytes)) / 10);
  const formattedSize = Number(BigInt(bytes) / BigInt(Math.pow(1024, i)));
  return `${formattedSize.toFixed(2)} ${sizes[i]}`;
}

export function formatNumber(number: number): string {
  const [integerPart, decimalPart] = number.toString().split(".");

  let result = "";
  const length = integerPart.length;

  for (let i = 0; i < length; i++) {
    if (i > 0 && (length - i) % 3 === 0) {
      result += "'";
    }
    result += integerPart[i];
  }

  if (decimalPart) {
    return result + "." + decimalPart;
  }

  return result;
}
