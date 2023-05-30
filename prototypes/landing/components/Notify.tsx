import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

export function Notify() {
  return (
    <div className="flex justify-center w-full max-w-sm space-x-2">
      <Input
        type="email"
        placeholder="Your email here.."
        className="text-white border-gray-500"
      />
      <Button type="submit" className="bg-orange-100">
        Notify
      </Button>
    </div>
  );
}
