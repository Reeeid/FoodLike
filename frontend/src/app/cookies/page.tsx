import type { Metadata } from "next";
import Link from "next/link";

export const metadata: Metadata = {
  title: "Cookieポリシー | FoodLike",
};

export default function CookiesPage() {
  return (
    <article className="card legal">
      <h1>Cookieポリシー</h1>
      <p className="legal-meta">制定日: 2026年11月10日</p>

      <p>
        FoodLike(以下「本サービス」)は、広告やアクセス解析のためのCookieを使用していません。本サービスがブラウザに保存するのは、サービスを動かすのに必要な次の情報だけです。
      </p>

      <h2>1. ブラウザに保存する情報</h2>
      <table>
        <thead>
          <tr>
            <th>名前</th>
            <th>保存先</th>
            <th>内容・目的</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>firebaseLocalStorageDb</td>
            <td>IndexedDB</td>
            <td>
              ログイン状態(Firebase
              Authenticationのトークン)。次回アクセス時にログインし直さずに済むようにします。
            </td>
          </tr>
          <tr>
            <td>foodlike:member</td>
            <td>localStorage</td>
            <td>ゲスト利用時のメンバー情報(ID・表示名)</td>
          </tr>
          <tr>
            <td>foodlike:lastRead:(グループID)</td>
            <td>localStorage</td>
            <td>グループごとに最後に読んだメッセージの番号。未読バッジの表示に使います。</td>
          </tr>
        </tbody>
      </table>
      <p>
        これらはお使いのブラウザの中にだけ保存され、未読位置はサーバーに送信しません。ログイン状態は、APIを利用するときの本人確認のためにサーバーへ送信します。
      </p>

      <h2>2. 外部サービスによるもの</h2>
      <p>
        Googleアカウントでログインする場合、ログイン画面の表示中にGoogleがCookieを使用することがあります。その取り扱いは
        <a
          href="https://policies.google.com/technologies/cookies?hl=ja"
          target="_blank"
          rel="noopener noreferrer"
        >
          Googleのポリシー
        </a>
        に従います。
      </p>

      <h2>3. 削除する方法</h2>
      <p>
        ログアウトするか、ブラウザの設定からこのサイトのデータ(Cookieとサイトデータ)を削除すると消えます。削除すると再ログインが必要になり、ゲストアカウントは元に戻せなくなります。
      </p>

      <h2>4. お問い合わせ</h2>
      <p>
        取り扱う情報全体については<Link href="/privacy">プライバシーポリシー</Link>
        をご覧ください。
      </p>
    </article>
  );
}
