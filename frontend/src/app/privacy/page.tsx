import type { Metadata } from "next";
import Link from "next/link";

export const metadata: Metadata = {
  title: "プライバシーポリシー | FoodLike",
};

const GITHUB_URL = "https://github.com/reeeid";

export default function PrivacyPage() {
  return (
    <article className="card legal">
      <h1>プライバシーポリシー</h1>
      <p className="legal-meta">制定日: 2026年11月10日</p>

      <p>
        FoodLike(以下「本サービス」)の運営者は、利用者の情報を以下のとおり取り扱います。
      </p>

      <h2>1. 取得する情報</h2>
      <table>
        <thead>
          <tr>
            <th>情報</th>
            <th>内容</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>アカウント情報</td>
            <td>
              メールアドレスまたはGoogleアカウントの識別子(Firebase
              Authenticationが管理)、表示名
            </td>
          </tr>
          <tr>
            <td>苦手なもの</td>
            <td>登録したジャンル・食材</td>
          </tr>
          <tr>
            <td>グループ情報</td>
            <td>参加しているグループ、招待コード</td>
          </tr>
          <tr>
            <td>メッセージ</td>
            <td>グループチャットの投稿、AI検索の質問とAIの回答、投稿日時</td>
          </tr>
          <tr>
            <td>アクセスログ</td>
            <td>
              アクセス日時、IPアドレス、リクエストしたパス(URLのクエリ文字列とメッセージ本文は記録しません)
            </td>
          </tr>
        </tbody>
      </table>

      <h2>2. 利用目的</h2>
      <ul>
        <li>本人確認とログイン状態の維持</li>
        <li>グループチャットとお店の提案機能の提供</li>
        <li>障害の調査、不正利用の防止</li>
      </ul>
      <p>
        広告配信、行動の分析(アクセス解析)、マーケティングには利用しません。
      </p>

      <h2>3. 苦手なものの取り扱い</h2>
      <p>
        登録した苦手なものは本人だけが閲覧・編集できます。お店の提案やAI検索では、グループ全員分をまとめて「誰のものか分からない」形にしてから使い、他のメンバーには誰が何を苦手かを表示しません。AIの回答に具体的な品目名が含まれた場合も、表示前に伏せ字にします。
      </p>

      <h2>4. メッセージの取り扱い(通信の秘密)</h2>
      <ul>
        <li>メッセージは、同じグループのメンバーにだけ表示されます。</li>
        <li>
          運営者は、法令に基づく場合、利用者の同意がある場合、またはサービスの維持・障害対応に必要な最小限の範囲を除き、メッセージの内容を閲覧・利用・第三者に提供しません。
        </li>
        <li>メッセージの本文はサーバーのログに記録しません。</li>
      </ul>

      <h2>5. 外部サービスへの送信</h2>
      <p>
        本サービスは、機能の提供のため次の外部サービスを利用しており、表の情報がそれぞれの事業者に送信されます。
      </p>
      <table>
        <thead>
          <tr>
            <th>送信先</th>
            <th>送信される情報</th>
            <th>目的</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>Google LLC(Firebase Authentication)</td>
            <td>メールアドレス、パスワード、Googleアカウント情報</td>
            <td>ログイン</td>
          </tr>
          <tr>
            <td>Google LLC(Gemini API)</td>
            <td>
              AI検索の質問文、グループ全員分をまとめた苦手なもの(個人とは紐づけません)
            </td>
            <td>AI回答の生成</td>
          </tr>
          <tr>
            <td>Tavily</td>
            <td>AI検索の質問文</td>
            <td>実在するお店の検索</td>
          </tr>
          <tr>
            <td>株式会社リクルート(ホットペッパー グルメサーチAPI)</td>
            <td>検索するエリア・予算(個人を特定する情報は送りません)</td>
            <td>お店の提案</td>
          </tr>
        </tbody>
      </table>
      <p>
        Gemini APIを無料枠で利用しているため、送信した質問文がGoogleのサービス改善に使われる場合があります(
        <a
          href="https://ai.google.dev/gemini-api/terms"
          target="_blank"
          rel="noopener noreferrer"
        >
          Gemini API 利用規約
        </a>
        )。AI検索には、個人情報や他人に知られたくない内容を入力しないでください。メンバー同士のチャット(「送信」ボタン)の内容は外部AIに送信しません。
      </p>

      <h2>6. 保存先</h2>
      <p>
        情報は次のクラウドサービス上に保存・処理されます。これらは国外のサーバーを含む場合があります。
      </p>
      <ul>
        <li>Vercel Inc.(画面の配信)</li>
        <li>Render Services, Inc.(APIサーバー)</li>
        <li>PingCAP(TiDB Cloud、データベース)</li>
        <li>Google LLC(Firebase、認証)</li>
      </ul>

      <h2>7. 第三者への提供</h2>
      <p>
        上記5・6を除き、法令に基づく場合を除いて、利用者の同意なく個人情報を第三者に提供しません。
      </p>

      <h2>8. 保存期間と削除</h2>
      <p>
        情報は本サービスの提供期間中保存します。アカウントやメッセージの削除を希望する場合は、下記の連絡先からご依頼ください。本人確認のうえ削除します。本サービスを終了する場合は、保存している情報を削除します。
      </p>

      <h2>9. ブラウザへの保存</h2>
      <p>
        ログイン状態などをブラウザに保存しています。詳しくは
        <Link href="/cookies">Cookieポリシー</Link>をご覧ください。
      </p>

      <h2>10. 改定</h2>
      <p>
        本ポリシーは必要に応じて改定し、本ページに掲載した時点から効力を持ちます。
      </p>

      <h2>お問い合わせ</h2>
      <p>
        情報の開示・訂正・削除の依頼やお問い合わせは、運営者の
        <a href={GITHUB_URL} target="_blank" rel="noopener noreferrer">
          GitHub
        </a>
        からご連絡ください。
      </p>
    </article>
  );
}
