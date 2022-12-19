import { useCallback, useState } from "react";
import { withEventValue } from "../utils";
import "./passwordModal.css"

type PasswordModalProps = {
    onSubmit: (lowPassword: string) => void
}
export const PasswordModal = (props: PasswordModalProps): JSX.Element => {
    const [hidden, setHidden] = useState(true)
    const [value, setValue] = useState("")
    const onSubmit = useCallback((e: React.FormEvent) => {
        e.preventDefault()
        props.onSubmit(value)
    }, [value])
    return (
        <div>
            <button className="password-modal-button">
                <span className="material-symbols-outlined" onClick={() => {
                    setHidden(false)
                }}>key</span>
            </button>
            <div className="password-modal" hidden={hidden}>
                <div className="background" onClick={() => setHidden(true)} />
                <div className="container">
                    <h1>パスワードを設定する</h1>
                    <p>データは暗号化され、 共通のパスワードを知っている人のみが 利用できるようになります。</p>
                    <form onSubmit={onSubmit}>
                        <input type="password" placeholder="******" value={value} onChange={withEventValue(setValue)} />
                        <button type="submit">
                            LOCK<span className="material-symbols-outlined">lock</span>
                        </button>
                    </form>
                    <p className="notice">※新たにパスワードを設定したチャットルームは全く別のものとして扱われるため、データは引き継がれません。</p>
                    <p className="notice">※パスワードが間違っていた際、注意書きは表示されません。あるはずのチャットが見つからないときは、誤ったパスワードを入力している可能性があります。</p>
                    <p className="notice">※パスワード管理アプリによって自動生成されたパスワードは比較的信頼性が高いといえます。ぜひご検討ください。</p>
                </div>
            </div>
        </div>
    )
}