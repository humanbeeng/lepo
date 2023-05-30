import { MainNavigation } from '@/components/MainNavigation';
import { Notify } from '@/components/Notify';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
} from '@/components/ui/card';

export default function Home() {
  return (
    <main className="bg-orange-100 min-h-screen">
      <MainNavigation />

      <section className="grid grid-cols-1 py-28  lg:py-60  items-center">
        <h1 className="font-heading px-2 justify-center  text-xl lg:text-5xl text-center lg:px-32">
          Lepo is an AI assistant who can help you
          <span className="underline">
            {' '}
            write, understand and navigate
          </span>{' '}
          code right inside your editor.
        </h1>
      </section>

      <section className="bg-black ">
        <div className="grid mx-2 justify-center justify-items-center gap-4 py-20 ">
          <h1 className="font-heading text-center  via-gray-300 text-gray-500 text-xl lg:text-3xl">
            Be the first to know when we launch..
          </h1>
          <Notify />
        </div>
      </section>

      <section className="bg-black pb-32">
        <div className="grid grid-cols-1 mx-2 lg:grid-cols-3 gap-4 lg:mx-10">
          <Card className="bg-transparent border-gray-500">
            <CardHeader className="font-extrabold text-orange-50 text-2xl">
              <p
                className="bg-clip-text
                text-transparent
                bg-gradient-to-r
                from-white
                via-orange-100
                to-orange-300"
              >
                Code Generation
              </p>
              <CardDescription className="text-gray-400 font-light pt-2">
                Examples
              </CardDescription>
            </CardHeader>

            <CardContent className="text-orange-50">
              <ul className="ml-6 list-disc [&>li]:mt-3">
                <li>
                  Add Hash-based URL load balancing to existing getNextBackend
                  method.
                </li>
                <li>
                  Write a new security filter which authorizes users based on
                  Authorization header.
                </li>
                <li>
                  Write a DB querier interface which can CRUD user struct.
                </li>
              </ul>
            </CardContent>
          </Card>
          <Card className="bg-transparent border-gray-500">
            <CardHeader className="font-extrabold text-orange-50 text-2xl">
              <p
                className="bg-clip-text
                text-transparent
                bg-gradient-to-r
                from-white
                via-orange-100
                to-orange-300"
              >
                Code Explanation
              </p>
              <CardDescription className="text-gray-400 font-light pt-2">
                Examples
              </CardDescription>
            </CardHeader>
            <CardContent className="text-yellow-50">
              <ul className="ml-6 list-disc [&>li]:mt-3">
                <li>Explain how user authentication works.</li>
                <li>
                  What happens if user tries to clear a cart which has no items
                  in it.
                </li>
                <li>Write a DB querier interface which can CRUD user struct</li>
              </ul>
            </CardContent>
          </Card>
          <Card className="bg-transparent border-gray-500">
            <CardHeader className="font-extrabold text-yellow-50 text-2xl">
              <p
                className="bg-clip-text
                text-transparent
                bg-gradient-to-r
                from-white
                via-orange-100
                to-orange-300"
              >
                Code Search
              </p>
              <CardDescription className="text-gray-400 font-light pt-2">
                Examples
              </CardDescription>
            </CardHeader>
            <CardContent className="text-yellow-50">
              <ul className="ml-6 list-disc [&>li]:mt-3">
                <li>
                  Where do we check if the user has permission to access payment
                  API?
                </li>
                <li>
                  Write a DB querier interface which can CRUD user struct.
                </li>
                <li>Do we have XML report renderer?</li>
              </ul>
            </CardContent>
          </Card>
        </div>
      </section>
    </main>
  );
}
