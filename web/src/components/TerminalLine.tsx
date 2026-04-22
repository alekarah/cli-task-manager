import type { TerminalLine as TLine } from '../types'
import styles from './Terminal.module.css'

interface Props {
  line: TLine
}

export default function TerminalLine({ line }: Props) {
  return <div className={`${styles.line} ${styles[line.type]}`}>{line.content}</div>
}
